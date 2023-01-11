package documents

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/luminosita/docrepo-bee/internal/interfaces/repositories/documents"
	documents2 "github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
	"github.com/luminosita/honeycomb/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPutDocumentGoodRequest(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	d := NewPutDocument(m.pdr)

	buffer := utils.NewNopWriteCloser(bytes.NewBuffer(make([]byte, 3*1024)))

	repoRequest := &documents.PutDocumentRepositorerRequest{Name: "laza", Size: 1234}
	repoResponse := &documents.PutDocumentRepositorerResponse{DocumentId: "pera", Writer: buffer}

	m.pdr.EXPECT().PutDocument(gomock.Eq(repoRequest)).Return(repoResponse, nil)

	res, err := d.Execute(&documents2.PutDocumenterRequest{Name: "laza", Size: 1234})

	assert.Nil(t, err)
	assert.Equal(t, "pera", res.DocumentId)
	assert.Equal(t, buffer, res.Writer)
}

func TestPutDocumentBadRequest(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	d := NewPutDocument(m.pdr)

	_, err := d.Execute(&documents2.PutDocumenterRequest{})

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Bad request")
}

func TestPutDocumentBadRequest2(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	d := NewPutDocument(m.pdr)

	_, err := d.Execute(nil)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Bad request")
}

func TestNewPutDocument(t *testing.T) {
	m := newMock(t)

	d := NewPutDocument(m.pdr)

	assert.Equal(t, &PutDocument{
		repo: m.pdr,
	}, d)
}
