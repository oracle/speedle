package oracle.speedle.ads.client;

import java.util.Map;
import java.util.logging.Level;
import java.util.logging.Logger;

import oracle.speedle.ads.client.rest.RestClient;

/**
 * The Speedle Client Factory. Users must use this Factory to construct a client instance.
 * @author jiefu
 *
 */
public class Factory {
    private static final String CLASS_NAME = Factory.class.getName();
    private static final Logger LOGGER = Logger.getLogger(CLASS_NAME);
    private static final Factory FACTORY = new Factory();

    private Factory() {}

    public static Factory getFactory() { return FACTORY; }

    /**
     * 
     * @param clientProperties
     * @return
     * @throws ADSException
     */
    public Client newClient(Map<String, Object> clientProperties) throws ADSException {
        final String methodName = "newClient";
        if (LOGGER.isLoggable(Level.FINER)) {
            LOGGER.entering(CLASS_NAME, methodName, clientProperties);
        }
        if (clientProperties == null) {
            String msg = "Argument clientProperties should not be null.";
            LOGGER.logp(Level.SEVERE, CLASS_NAME, methodName, msg);
            throw new IllegalArgumentException(msg);
        }

        Object restHost = clientProperties.get(Client.ClientPropertyName.REST_ENDPOINT.name());

        if (LOGGER.isLoggable(Level.FINE)) {
            LOGGER.logp(Level.FINE, CLASS_NAME, methodName, String.format("Rest Host: %s", restHost));
        }

        if (restHost == null ) {
            String msg = String.format("Properties of %s should not be null.", Client.ClientPropertyName.REST_ENDPOINT);
            LOGGER.logp(Level.SEVERE, CLASS_NAME, methodName, msg);
            throw new IllegalArgumentException(msg);
        }
       
        if (!(restHost instanceof String)) {
            String msg = "Property " + Client.ClientPropertyName.REST_ENDPOINT + " should be an object of String";
            LOGGER.logp(Level.SEVERE, CLASS_NAME, methodName, msg);
            throw new IllegalArgumentException(msg);
        }

        Client client = new RestClient((String)restHost);
        if (LOGGER.isLoggable(Level.FINER)) {
            LOGGER.exiting(CLASS_NAME, methodName, client);
        }
        return client;
    }
}
