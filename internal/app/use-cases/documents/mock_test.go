package documents

import (
	"github.com/golang/mock/gomock"
	"github.com/luminosita/docrepo-bee/internal/interfaces/repositories/documents/mocks"
	"testing"
)

type Mock struct {
	t    *testing.T
	ctrl *gomock.Controller

	pdr  *mocks.MockPutDocumentRepositorer
	gdr  *mocks.MockGetDocumentRepositorer
	gidr *mocks.MockGetDocumentInfoRepositorer
}

func newMock(t *testing.T) (m *Mock) {
	m = &Mock{}
	m.t = t
	m.ctrl = gomock.NewController(t)

	m.pdr = mocks.NewMockPutDocumentRepositorer(m.ctrl)
	m.gdr = mocks.NewMockGetDocumentRepositorer(m.ctrl)
	m.gidr = mocks.NewMockGetDocumentInfoRepositorer(m.ctrl)

	return
}

func setupTest(m *Mock) func() {
	if m == nil {
		panic("Mock not initialized")
	}

	return func() {
	}
}
