# Speedle Embedeed Sample

# Build
Please make sure envionment variables GOROOT and GOPATH have been properly set. And you are currently under this folder.
```bash
$ go build -o expenses_sample cmd/expenses_sample
```

After building, a executable file `expenses_sample` can be found under this folder.

# Policies for this sample
The policy definitations for this sample can be found in file expenses.spdl
```
[service.expenses]
[policy]
GRANT ROLE employee get, post, delete /reports
GRANT ROLE auditor get, modify /reports
[rolepolicy]
GRANT USER alice employee
```

In this file, there is a service named `expenses`, and under this service, there are two policies defined, role `employee` can `get`,
`post` and `delete` resource `/reports`; and role `auditor` can `get` and `modify` resource `/reports`.

And user `alice` is an employee.

# Run and test the sample
```bash
# Run the sample, the sample will listen on localhost:8080
$ ./expenses_sample expenses.spdl
```

```bash
# Test if user alice has permission to delete /reports
$ curl -X DELETE -u alice:afdsa http://localhost:8080/reports
deleteing an expense report is done

# Test if user bob has permission to delete /reports
$ curl -X DELETE -u bob:afdsa http://localhost:8080/reports
forbidden.
```

# Add bob as an employee, and test
```bash
# Add user bob as an employee
$ echo "GRANT USER bob employee" >> expenses.spdl
```

```bash
# Test if bob has permission to delete /reports
$ curl -X DELETE -u bob:afdsa http://localhost:8080/reports
deleteing an expense report is done
```
