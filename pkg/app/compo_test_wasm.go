package app

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type boo struct {
	Compo

	Value int

	onDismount func()
	onUpdate   func()
}

func (b *boo) OnDismount() {
	if b.onDismount != nil {
		b.onDismount()
	}
}

func (b *boo) OnUpdate() {
	if b.onUpdate != nil {
		b.onUpdate()
	}
}

func (b *boo) Render() UI {
	return Text("foo")
}

type booWithDefaultRender struct {
	Compo
}

func TestCompoUnmountedUpdate(t *testing.T) {
	tests := []struct {
		scenario string
		compo    Composer
	}{
		{
			scenario: "component with redefined render is updated",
			compo:    &boo{},
		},
		{
			scenario: "component without redefined render is updated",
			compo:    &booWithDefaultRender{},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			dispatcher = func(f func()) {
				f()
			}
			defer func() {
				dispatcher = Dispatch
			}()

			test.compo.Update()
		})
	}
}

func TestCompoDismount(t *testing.T) {
	called := false

	c := &boo{
		onDismount: func() {
			called = true
		},
	}

	mount(c)
	c.dismount()
	require.True(t, called)
}

type navTest struct {
	Compo

	subcompo UI
	onNav    func(*url.URL)
}

func (n *navTest) OnNav(u *url.URL) {
	if n.onNav != nil {
		n.onNav(u)
	}
}

func (n *navTest) Render() UI {
	return Div().Body(
		n.subcompo,
	)
}

func TestNavigator(t *testing.T) {
	bcalled := false
	b := &navTest{
		onNav: func(u *url.URL) {
			bcalled = true
		},
	}

	acalled := false
	a := &navTest{
		subcompo: b,
		onNav: func(u *url.URL) {
			acalled = true
		},
	}

	err := mount(a)
	require.NoError(t, err)

	require.False(t, acalled)
	require.False(t, bcalled)

	triggerOnNav(a, nil)
	require.True(t, true)
	require.True(t, true)
}
