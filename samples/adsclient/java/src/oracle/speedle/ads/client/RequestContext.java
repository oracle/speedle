package oracle.speedle.ads.client;

import java.util.Map;

public class RequestContext {
    private Subject subject = null;
    private String serviceName = null;
    private String resource = null;
    private String action = null;
    private Map<String, Object> attributes = null;
    private Map<String, String> tokens = null;

    public RequestContext() { }

    public void setSubject(Subject subject) {
        this.subject = subject;
    }
    public Subject getSubject() { return subject; }


    public void setServiceName(String serviceName) {
        this.serviceName = serviceName;
    }
    public String getServiceName() { return serviceName; }

    public void setResource(String resource) {
        this.resource = resource;
    }
    public String getResource() { return resource; }

    public void setAction(String action) {
        this.action = action;
    }
    public String getAction() { return action; }

    public void setAttributes(Map<String, Object> attributes) {
        this.attributes = attributes;
    }
    public Map<String, Object> getAttributes() { return attributes; }

    public void setTokens(Map<String, String> tokens) {
        this.tokens = tokens;
    }
    public Map<String, String> getTokens() { return tokens; }
}
