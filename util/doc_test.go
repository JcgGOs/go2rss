package util

import (
	"testing"
)

func TestAbsHref(t *testing.T) {

	tests := []struct {
		name string
		link string
		href string
		want string
	}{
		{name: "1", want: "http://company.com/access", link: "http://company.com/blog", href: "http://company.com/access"},
		{name: "2", want: "http://company.com/access", link: "http://company.com/blog", href: "//company.com/access"},
		{name: "3", want: "http://company.com/access", link: "http://company.com/blog", href: "access"},
		{name: "4", want: "http://company.com/blog/access", link: "http://company.com/blog/", href: "access"},
		{name: "5", want: "http://company.com/blog/../access", link: "http://company.com/blog/", href: "../access"},
		{name: "6", want: "http://company.com/access", link: "http://company.com/blog/", href: "/access"},
		{name: "7", want: "http://company.com/access", link: "http://company2.com/blog/", href: "http://company.com/access"},
		{name: "8", want: "http://company.com/access", link: "http://company.com", href: "/access"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsHref(tt.href, tt.link); got != tt.want {
				t.Errorf("AbsHref() = %v, want %v", got, tt.want)
			}
		})
	}
}
