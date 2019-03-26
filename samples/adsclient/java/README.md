# Speedle Java Client

This is a Java Client for Speedle ADS to simplify ADS calls.

# Interface
```java
package oracle.speedle.ads.client;
public interface Client {
    public enum ClientPropertyName {
        /**
         * 
         * The property value should be a string object. Speedle ADS RESTful service host name.
         */
        REST_ENDPOINT,
    }

    /**
     * isAllowed evaluates policies with inputed request context, and returns the decision.
     * @param context The request context.
     * @return The decision
     * @throws ADSException Error happened. For example, invalid request context, network issue, server internal issue, etc. 
     */
    public boolean isAllowed(RequestContext context) throws ADSException;
}
```

# How to construct a REST client instance
```java
        Map<String, Object> properties = new HashMap<String, Object>();
        properties.put(Client.ClientPropertyName.REST_ENDPOINT.name(), "localhost:6734");    // Set REST Endpoint to localhost:6734
        Client client = Factory.getFactory().newClient(properties);    // Construct a client instance with properties
```

# Example
The file [Example.java](src/example/Example.java) showes how to call ADS promatically.

# How to compile this client
## Require
* JDK 1.8 has been installed, and related environment variables JAVA_HOME and PATH were set
* Maven has been installed.

## Command to build client classes
```bash
# Enter the java client home
$ cd samples/adsclient/java
# Compile all classes
$ mvn compile
```

# Q&A
## How to connect ADS with this client behind a HTTP proxy.
Run JVM with three system properties:
* https.proxyHost: The proxy IP Address or host name for HTTPS connections
* https.proxyPort: The proxy port for HTTPS connections
* http.nonProxyHosts: Hosts that should not be connected with the proxy.

For more details, access [Java Networking and Proxies](https://docs.oracle.com/javase/8/docs/technotes/guides/net/proxies.html)

## JDK Version
JDK 1.8
