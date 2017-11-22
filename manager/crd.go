package manager

import (
	"fmt"
	"reflect"

	"strings"

	"github.com/pkg/errors"
	catalogv1 "github.com/rancher/types/apis/catalog.cattle.io/v1"
	"github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CatalogNameLabel = "io.cattle.catalog.name"
)

// update will sync templates with catalog without costing too much
func (m *manager) update(catalog *catalogv1.Catalog, templates []catalogv1.Template) error {
	logrus.Debugf("Syncing catalog %s with templates", catalog.Name)
	existingTemplates, err := m.templateClient.List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", CatalogNameLabel, catalog.Name),
	})
	if err != nil {
		return err
	}

	templatesByName := map[string]catalogv1.Template{}
	for _, template := range templates {
		if template.Spec.FolderName == "" {
			continue
		} else if template.Spec.Base == "" && template.Spec.FolderName != "" {
			template.Name = fmt.Sprintf("%s-%s", catalog.Name, template.Spec.FolderName)
		} else {
			template.Name = fmt.Sprintf("%s-%s-%s", catalog.Name, template.Spec.Base, template.Spec.FolderName)
		}
		template.Name = strings.ToLower(template.Name)
		templatesByName[template.Name] = template
	}

	existingTemplatesByName := map[string]catalogv1.Template{}
	for _, template := range existingTemplates.Items {
		existingTemplatesByName[template.Name] = template
	}

	// templates is the one we should update, so for all the templates that were in existingTemplates
	// 1. if it doesn't exist in templates, delete them
	// 2. if it exists but has changed, update it
	// 3. if it exists but not changed, keep it unmodified
	for name, existingTemplate := range existingTemplatesByName {
		template, ok := templatesByName[name]
		if !ok {
			// delete the template
			logrus.Debugf("Deleting templates %s", name)
			if err := m.templateClient.Delete(name, &metav1.DeleteOptions{}); err != nil {
				return errors.Wrapf(err, "failed to delete template %s", template.Name)
			}
		}

		if !reflect.DeepEqual(template.Spec, existingTemplate.Spec) {
			updateTemplate, err := m.templateClient.Get(name, metav1.GetOptions{})
			if err != nil && !kerrors.IsNotFound(err) {
				return err
			} else if kerrors.IsNotFound(err) {
				continue
			}
			updateTemplate.Spec = template.Spec
			logrus.Debugf("Updating template %s", name)
			_, err = m.templateClient.Update(updateTemplate)
			if err != nil {
				if strings.Contains(err.Error(), "request is too large") || strings.Contains(err.Error(), "exceeding the max size") {
					updateTemplate.Spec.Icon = ""
					if _, err := m.templateClient.Update(updateTemplate); err != nil {
						return err
					}
					return nil
				}
				return errors.Wrapf(err, "failed to update template %s", template.Name)
			}
		}
	}

	// for templates that exist in template but not in existingTemplates, we should create them
	for name, template := range templatesByName {
		if _, ok := existingTemplatesByName[name]; !ok {
			template.OwnerReferences = []metav1.OwnerReference{
				{
					APIVersion: catalog.APIVersion,
					Kind:       catalog.Kind,
					Name:       catalog.Name,
					UID:        catalog.UID,
					Controller: &[]bool{true}[0],
				},
			}
			template.Kind = catalogv1.TemplateGroupVersionKind.Kind
			template.APIVersion = catalogv1.TemplateGroupVersionKind.Group + "/" + catalogv1.TemplateGroupVersionKind.Version
			template.Labels = map[string]string{}
			template.Labels[CatalogNameLabel] = catalog.Name
			logrus.Debugf("Creating template %s", template.Name)
			_, err := m.templateClient.Create(&template)
			if err != nil {
				// hack for the image size that are too big
				if strings.Contains(err.Error(), "request is too large") || strings.Contains(err.Error(), "exceeding the max size") {
					template.Spec.Icon = ""
					if _, err := m.templateClient.Create(&template); err != nil {
						return err
					}
				}
				return err
			}
		}
	}
	return nil
}
