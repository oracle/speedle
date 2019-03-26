package example;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

import oracle.speedle.ads.client.Client;
import oracle.speedle.ads.client.Factory;
import oracle.speedle.ads.client.Principal;
import oracle.speedle.ads.client.Principal.PrincipalType;
import oracle.speedle.ads.client.RequestContext;
import oracle.speedle.ads.client.Subject;

public class Example {
    public static void main(String args[]) throws Exception {
        if (args.length != 2) {
            System.err.printf("Arguments: <PMS Endpoint> <Speedle Service>\n\nPMS Endpoint:    localhost:6734 for example\nSpeedle Service: testsvc for example\n\n");
            System.exit(1);
        }

        Map<String, Object> properties = new HashMap<String, Object>();
        properties.put(Client.ClientPropertyName.REST_ENDPOINT.name(), args[0]);

        Client client = Factory.getFactory().newClient(properties);

        Subject subject = new Subject();
        Principal princ = new Principal();
        princ.setType(PrincipalType.USER);
        princ.setName("testuser");
        subject.setPrincipals(Arrays.asList(princ));
        RequestContext context = new RequestContext();
        context.setSubject(subject);
        context.setServiceName(args[1]);
        context.setResource("testres");
        context.setAction("read");

        System.out.println("Evaluation result is: " + client.isAllowed(context));
    }
}
