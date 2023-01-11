package documents

import (
	"github.com/golang/mock/gomock"
	"github.com/luminosita/docrepo-bee/internal/interfaces/repositories/documents"
	documents2 "github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestGetDocumentGoodRequest(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	d := NewGetDocument(m.gdr)

	reader := io.NopCloser(strings.NewReader("test"))

	repoRequest := &documents.GetDocumentRepositorerRequest{DocumentId: "pera"}
	repoResponse := &documents.GetDocumentRepositorerResponse{Name: "laza", Size: 1234, Reader: reader}

	m.gdr.EXPECT().GetDocument(gomock.Eq(repoRequest)).Return(repoResponse, nil)

	res, err := d.Execute(&documents2.GetDocumenterRequest{DocumentId: "pera"})

	assert.Nil(t, err)
	assert.Equal(t, reader, res.Reader)
	assert.Equal(t, int64(1234), res.Size)
	assert.Equal(t, "laza", res.Name)
}

func TestGetDocumentBadRequest(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	d := NewGetDocument(m.gdr)

	_, err := d.Execute(&documents2.GetDocumenterRequest{})

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Bad request")
}

func TestGetDocumentBadRequest2(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	d := NewGetDocument(m.gdr)

	_, err := d.Execute(nil)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Bad request")
}

func TestNewGetDocument(t *testing.T) {
	m := newMock(t)

	d := NewGetDocument(m.gdr)

	assert.Equal(t, &GetDocument{
		repo: m.gdr,
	}, d)
}
