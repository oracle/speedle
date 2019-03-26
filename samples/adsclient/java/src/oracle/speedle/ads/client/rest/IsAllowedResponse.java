package oracle.speedle.ads.client.rest;

class IsAllowedResponse {
    private boolean allowed = false;
    private String message = null;

    public void setAllowed(boolean allowed) { this.allowed = allowed; }
    public boolean isAllowed() { return allowed; }
    
    public void setMessage(String message) { this.message = message; }
    public String getMessage() { return message; }
}
