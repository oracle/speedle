//Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package benchmarks

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/oracle/speedle/api/ads"
	"github.com/oracle/speedle/pkg/eval"
	_ "github.com/oracle/speedle/pkg/store/file"
)

func BenchmarkOne(b *testing.B) {
	filePath, err := writePolicyFile(b.Name(), 1, 1, simplePolicyWriter, simpleRolePolicyWriter)
	if err != nil {
		b.Fatalf("unable to write to policy file because of error %s", err)
	}
	defer os.Remove(filePath)

	rc := ads.RequestContext{
		Subject: &ads.Subject{
			Principals: []*ads.Principal{
				{
					Type: "user",
					Name: "user1-1",
				},
			},
		},
		ServiceName: "bench",
		Action:      "read",
		Resource:    "/books/book1",
	}

	runTest(b, rc, filePath, true)
}

func BenchmarkTiny(b *testing.B) {
	filePath, err := writePolicyFile(b.Name(), 10, 10, simplePolicyWriter, simpleRolePolicyWriter)
	if err != nil {
		b.Fatalf("unable to write to policy file because of error %s", err)
	}
	defer os.Remove(filePath)

	rc := ads.RequestContext{
		Subject: &ads.Subject{
			Principals: []*ads.Principal{
				{
					Type: "user",
					Name: "user5-7",
				},
			},
		},
		ServiceName: "bench",
		Action:      "read",
		Resource:    "/books/book5",
	}

	runTest(b, rc, filePath, true)
}

func BenchmarkSmall(b *testing.B) {
	filePath, err := writePolicyFile(b.Name(), 100, 10, simplePolicyWriter, simpleRolePolicyWriter)
	if err != nil {
		b.Fatalf("unable to write to policy file because of error %s", err)
	}
	defer os.Remove(filePath)

	rc := ads.RequestContext{
		Subject: &ads.Subject{
			Principals: []*ads.Principal{
				{
					Type: "user",
					Name: "user50-7",
				},
			},
		},
		ServiceName: "bench",
		Action:      "read",
		Resource:    "/books/book50",
	}

	runTest(b, rc, filePath, true)
}

func BenchmarkMedium(b *testing.B) {
	filePath, err := writePolicyFile(b.Name(), 1000, 10, simplePolicyWriter, simpleRolePolicyWriter)
	if err != nil {
		b.Fatalf("unable to write to policy file because of error %s", err)
	}
	defer os.Remove(filePath)

	rc := ads.RequestContext{
		Subject: &ads.Subject{
			Principals: []*ads.Principal{
				{
					Type: "user",
					Name: "user500-7",
				},
			},
		},
		ServiceName: "bench",
		Action:      "read",
		Resource:    "/books/book500",
	}

	runTest(b, rc, filePath, true)
}

func BenchmarkLarge(b *testing.B) {
	filePath, err := writePolicyFile(b.Name(), 10000, 10, simplePolicyWriter, simpleRolePolicyWriter)
	if err != nil {
		b.Fatalf("unable to write to policy file because of error %s", err)
	}
	defer os.Remove(filePath)

	rc := ads.RequestContext{
		Subject: &ads.Subject{
			Principals: []*ads.Principal{
				{
					Type: "user",
					Name: "user5000-7",
				},
			},
		},
		ServiceName: "bench",
		Action:      "read",
		Resource:    "/books/book5000",
	}

	runTest(b, rc, filePath, true)
}

func BenchmarkHuge(b *testing.B) {
	filePath, err := writePolicyFile(b.Name(), 100000, 10, simplePolicyWriter, simpleRolePolicyWriter)
	if err != nil {
		b.Fatalf("unable to write to policy file because of error %s", err)
	}
	defer os.Remove(filePath)

	rc := ads.RequestContext{
		Subject: &ads.Subject{
			Principals: []*ads.Principal{
				{
					Type: "user",
					Name: "user50000-7",
				},
			},
		},
		ServiceName: "bench",
		Action:      "read",
		Resource:    "/books/book50000",
	}

	runTest(b, rc, filePath, true)
}

func BenchmarkLargeExp(b *testing.B) {
	filePath, err := writePolicyFile(b.Name(), 10000, 10, func(w io.Writer, pno int) error {
		_, err := fmt.Fprintf(w, "GRANT ROLE role%d read expr:/books/type%d/.*\n", pno, pno)
		return err
	}, simpleRolePolicyWriter)
	if err != nil {
		b.Fatalf("unable to write to policy file because of error %s", err)
	}
	defer os.Remove(filePath)

	rc := ads.RequestContext{
		Subject: &ads.Subject{
			Principals: []*ads.Principal{
				{
					Type: "user",
					Name: "user5000-7",
				},
			},
		},
		ServiceName: "bench",
		Action:      "read",
		Resource:    "/books/type5000/fdsa;qlwejjmkfkld'sa",
	}

	runTest(b, rc, filePath, true)
}

func BenchmarkLargeCond(b *testing.B) {
	filePath, err := writePolicyFile(b.Name(), 10000, 10, func(w io.Writer, pno int) error {
		_, err := fmt.Fprintf(w, "GRANT ROLE role%d read /books/book%d if att1 == \"val1\" && att2 == \"val2\"\n", pno, pno)
		return err
	}, simpleRolePolicyWriter)
	if err != nil {
		b.Fatalf("unable to write to policy file because of error %s", err)
	}
	defer os.Remove(filePath)

	rc := ads.RequestContext{
		Subject: &ads.Subject{
			Principals: []*ads.Principal{
				{
					Type: "user",
					Name: "user5000-7",
				},
			},
		},
		ServiceName: "bench",
		Action:      "read",
		Resource:    "/books/book5000",
		Attributes: map[string]interface{}{
			"att1": "val1",
			"att2": "val2",
		},
	}

	runTest(b, rc, filePath, true)
}

func BenchmarkLargePerms(b *testing.B) {
	filePath, err := writePolicyFile(b.Name(), 10000, 10, func(w io.Writer, pno int) error {
		for i := 1; i <= 5; i++ {
			_, err := fmt.Fprintf(w, "GRANT ROLE role%d read /books/book%d\n", pno, i)
			if err != nil {
				return nil
			}
		}
		return nil
	}, simpleRolePolicyWriter)
	if err != nil {
		b.Fatalf("unable to write to policy file because of error %s", err)
	}
	defer os.Remove(filePath)

	ev, err := eval.NewFromFile(filePath, false)
	if err != nil {
		b.Fatalf("Unable to initialize evaluator due to error %s.", err)
	}

	rc := ads.RequestContext{
		Subject: &ads.Subject{
			Principals: []*ads.Principal{
				{
					Type: "user",
					Name: "user5000-7",
				},
			},
		},
		ServiceName: "bench",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		perms, err := ev.GetAllGrantedPermissions(rc)
		if err != nil {
			b.Fatalf("error occured %s.", err)
		}
		if len(perms) != 5 {
			b.Fatal("wrong returned permssions.", err)
		}
	}
	b.StopTimer()
}

func simplePolicyWriter(w io.Writer, pno int) error {
	_, err := fmt.Fprintf(w, "GRANT ROLE role%d read /books/book%d\n", pno, pno)
	return err
}

func simpleRolePolicyWriter(w io.Writer, pno, uc int) error {
	if _, err := fmt.Fprintf(w, "GRANT USER user%d-1", pno); err != nil {
		return err
	}
	for i := 2; i <= uc; i++ {
		if _, err := fmt.Fprintf(w, ", USER user%d-%d", pno, i); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w, " role%d\n", pno)
	return err
}

func writePolicyFile(fname string, pc, uc int, pw func(w io.Writer, pno int) error, rpw func(w io.Writer, pno, uc int) error) (string, error) {
	pf, err := ioutil.TempFile(os.TempDir(), fname+"_*.spdl")
	if err != nil {
		return "", err
	}
	defer pf.Close()

	if _, err := fmt.Fprintln(pf, "[service.bench]"); err != nil {
		return "", err
	}
	if _, err := fmt.Fprintln(pf, "[policy]"); err != nil {
		return "", err
	}
	for i := 1; i <= pc; i++ {
		if err := pw(pf, i); err != nil {
			return "", err
		}
	}
	if _, err := fmt.Fprintln(pf, "[rolepolicy]"); err != nil {
		return "", err
	}
	for i := 1; i <= pc; i++ {
		if err := rpw(pf, i, uc); err != nil {
			return "", err
		}
	}
	return pf.Name(), nil
}

func runTest(b *testing.B, rc ads.RequestContext, pfloc string, exp bool) {
	ev, err := eval.NewFromFile(pfloc, false)
	if err != nil {
		b.Fatalf("Unable to initialize evaluator due to error %s.", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		allowed, _, _ := ev.IsAllowed(rc)
		if !allowed == exp {
			b.Fatal("wrong decision result.")
		}
	}
	b.StopTimer()
}
