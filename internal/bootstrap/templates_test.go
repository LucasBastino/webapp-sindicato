package bootstrap

import (
	"testing"

	"github.com/LucasBastino/webapp-sindicato/internal/common/functiontemplates"
	"github.com/gofiber/template/html/v2"
)

func TestTemplatesLoad(t *testing.T) {
	engine := html.New("../views", ".html")
	engine.AddFunc("formatAmountAR", functiontemplates.FormatAmountAR)
	engine.AddFunc("formatAmountInput", functiontemplates.FormatAmountInput)

	if err := engine.Load(); err != nil {
		t.Fatalf("failed to load templates: %v", err)
	}
}
