package provider

import (
	"strings"
	"testing"

	"github.com/yetone/magpie/internal/catalog"
)

func TestQianfanTokenPlan(t *testing.T) {
	p, err := FromPreset("baidu-qianfan")
	if err != nil {
		t.Fatal(err)
	}
	// the plan's own endpoints, one base per protocol family; the plans
	// serve no model list, so the preset's models are the picker's list
	if p.Chat != "https://qianfan.baidubce.com/v2/tokenplan/personal" || p.Responses != p.Chat ||
		p.Anthropic != "https://qianfan.baidubce.com/anthropic/tokenplan/personal" {
		t.Fatalf("endpoints: %q %q %q", p.Chat, p.Responses, p.Anthropic)
	}
	// the personal and the enterprise plan and pay as you go, each plan's
	// bases under its own path, pay as you go at the v2 root
	for _, r := range Preset("baidu-qianfan").Regions {
		var chat, anthropic string
		switch r.ID {
		case "personal":
			chat, anthropic = "https://qianfan.baidubce.com/v2/tokenplan/personal", "https://qianfan.baidubce.com/anthropic/tokenplan/personal"
		case "team":
			chat, anthropic = "https://qianfan.baidubce.com/v2/tokenplan/team", "https://qianfan.baidubce.com/anthropic/tokenplan/team"
		case "api":
			chat, anthropic = "https://qianfan.baidubce.com/v2", "https://qianfan.baidubce.com/anthropic"
		default:
			t.Fatalf("region: %+v", r)
		}
		if r.Chat != chat || r.Responses != chat || r.Anthropic != anthropic {
			t.Fatalf("region %s: %q %q %q", r.ID, r.Chat, r.Responses, r.Anthropic)
		}
	}
	if len(Preset("baidu-qianfan").Regions) != 3 {
		t.Fatalf("regions: %d", len(Preset("baidu-qianfan").Regions))
	}
	for _, id := range Preset("baidu-qianfan").Models {
		if id != strings.ToLower(id) {
			t.Fatalf("model ids are lowercase on the plan: %q", id)
		}
	}
	// the plan's models, qianfan-code-latest following the console's pick first
	if got := p.planModels(nil); len(got) != len(Preset("baidu-qianfan").Models) || got[0].ID != "qianfan-code-latest" {
		t.Fatalf("plan's: %+v", got)
	}
	// a list the plan gave is kept whole
	if got := p.planModels([]catalog.Model{{ID: "glm-5.3"}}); len(got) != 1 || got[0].ID != "glm-5.3" {
		t.Fatalf("listed: %+v", got)
	}
	// an entry imported from another app at the plan's endpoints is the preset
	im, _ := imported("Qianfan", "bce-v3/x", endpoints{anthropic: "https://qianfan.baidubce.com/anthropic/tokenplan/personal"}, nil)
	if im.Preset != "baidu-qianfan" || im.Icon != "baiducloud-color" {
		t.Fatalf("imported: %+v", im)
	}
}
