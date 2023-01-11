package documents

import (
	"github.com/golang/mock/gomock"
	"github.com/luminosita/docrepo-bee/internal/interfaces/repositories/documents"
	documents2 "github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetDocumentInfoGoodRequest(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	d := NewGetDocumentInfo(m.gidr)

	ttt := time.Unix(100, 10)

	repoRequest := &documents.GetDocumentInfoRepositorerRequest{DocumentId: "pera"}
	repoResponse := &documents.GetDocumentInfoRepositorerResponse{Name: "laza", Size: 1234, UploadDate: ttt}

	m.gidr.EXPECT().GetDocumentInfo(gomock.Eq(repoRequest)).Return(repoResponse, nil)

	res, err := d.Execute(&documents2.GetDocumentInfoerRequest{DocumentId: "pera"})

	assert.Nil(t, err)
	assert.Equal(t, ttt, res.UploadDate)
	assert.Equal(t, int64(1234), res.Size)
	assert.Equal(t, "laza", res.Name)
}

func TestGetDocumentInfoBadRequest(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	d := NewGetDocumentInfo(m.gidr)

	_, err := d.Execute(&documents2.GetDocumentInfoerRequest{})

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Bad request")
}

func TestGetDocumentInfoBadRequest2(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	d := NewGetDocumentInfo(m.gidr)

	_, err := d.Execute(nil)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Bad request")
}

func TestNewGetDocumentInfo(t *testing.T) {
	m := newMock(t)

	d := NewGetDocumentInfo(m.gidr)

	assert.Equal(t, &GetDocumentInfo{
		repo: m.gidr,
	}, d)
}
