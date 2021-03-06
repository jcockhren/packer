package dockertag

import (
	"bytes"
	"github.com/mitchellh/packer/builder/docker"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/post-processor/docker-import"
	"testing"
)

func testConfig() map[string]interface{} {
	return map[string]interface{}{
		"repository": "foo",
		"tag":        "bar",
	}
}

func testPP(t *testing.T) *PostProcessor {
	var p PostProcessor
	if err := p.Configure(testConfig()); err != nil {
		t.Fatalf("err: %s", err)
	}

	return &p
}

func testUi() *packer.BasicUi {
	return &packer.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	}
}

func TestPostProcessor_ImplementsPostProcessor(t *testing.T) {
	var _ packer.PostProcessor = new(PostProcessor)
}

func TestPostProcessor_PostProcess(t *testing.T) {
	driver := &docker.MockDriver{}
	p := &PostProcessor{Driver: driver}
	_, err := common.DecodeConfig(&p.config, testConfig())
	if err != nil {
		t.Fatalf("err %s", err)
	}

	artifact := &packer.MockArtifact{
		BuilderIdValue: dockerimport.BuilderId,
		IdValue:        "1234567890abcdef",
	}

	result, keep, err := p.PostProcess(testUi(), artifact)
	if _, ok := result.(packer.Artifact); !ok {
		t.Fatal("should be instance of Artifact")
	}
	if !keep {
		t.Fatal("should keep")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !driver.TagImageCalled {
		t.Fatal("should call TagImage")
	}
	if driver.TagImageImageId != "1234567890abcdef" {
		t.Fatal("bad image id")
	}
	if driver.TagImageRepo != "foo:bar" {
		t.Fatal("bad repo")
	}
}
