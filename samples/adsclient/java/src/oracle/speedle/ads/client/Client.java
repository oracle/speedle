package oracle.speedle.ads.client;

/**
 * The interface of speedle Client.
 * @author jiefu
 *
 */
public interface Client {
    /**
     * The enum for all client properties names.
     * @author jiefu
     *
     */
    public enum ClientPropertyName {
        /**
         * 
         * The property value should be a string object. speedle ADS RESTful service host name.
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
