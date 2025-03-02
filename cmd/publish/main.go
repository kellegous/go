package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/pflag"
)

type Builder struct {
	Name string
}

func (b *Builder) Shutdown(ctx context.Context) error {
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"buildx",
		"rm",
		b.Name)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func StartBuilder(ctx context.Context) (*Builder, error) {
	var key [8]byte
	if _, err := rand.Read(key[:]); err != nil {
		return nil, err
	}
	name := fmt.Sprintf("go-build-%s", hex.EncodeToString(key[:]))

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"buildx",
		"create",
		"--name",
		name)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &Builder{Name: name}, nil
}

func (b *Builder) Build(
	ctx context.Context,
	image string,
	tags []string,
	platforms []string,
	buildArgs []string,
) error {
	args := []string{
		"buildx",
		"build",
		"--builder",
		b.Name,
		"--platform",
		strings.Join(platforms, ","),
	}

	for _, arg := range buildArgs {
		args = append(args, "--build-arg", arg)
	}

	for _, tag := range tags {
		args = append(args,
			"-t",
			fmt.Sprintf("%s:%s", image, tag))
	}

	args = append(args, "--push", ".")

	cmd := exec.CommandContext(
		ctx,
		"docker",
		args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func run(
	ctx context.Context,
	image string,
	tags []string,
	platforms []string,
	buildArgs []string,
) error {
	builder, err := StartBuilder(ctx)
	if err != nil {
		return err
	}
	defer builder.Shutdown(ctx)

	if err := builder.Build(ctx, image, tags, platforms, buildArgs); err != nil {
		return err
	}

	return nil
}

func main() {
	var tags []string
	var platforms []string
	var buildArgs []string
	var image string

	pflag.StringSliceVar(
		&tags,
		"tag",
		[]string{"latest"},
		"tags to apply to the image")
	pflag.StringSliceVar(
		&platforms,
		"platform",
		[]string{"linux/amd64", "linux/arm64"},
		"platforms to build for")
	pflag.StringSliceVar(
		&buildArgs,
		"build-arg",
		[]string{},
		"build arguments to pass to the image")
	pflag.StringVar(
		&image,
		"image",
		"kellegous/go",
		"the name of the image to build")
	pflag.Parse()

	if err := run(
		context.Background(),
		image, tags,
		platforms,
		buildArgs,
	); err != nil {
		log.Panic(err)
	}
}
