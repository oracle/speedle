package oracle.speedle.ads.client;

/**
 * Exception for all ADS calls
 * @author jiefu
 *
 */
public class ADSException extends Exception {
    private static final long serialVersionUID = -5945484282574506595L;

    /**
     * Default constructor to instance an ADSException object.
     */
    public ADSException() {
        super();
    }

    /**
     * Constructor to instance an ADSException object with an error message.
     * @param message The error message.
     */
    public ADSException(String message) {
        super(message);
    }

    /**
     * Constructor to instance an ADSException object with an error message and the cause.
     * @param message The error message.
     * @param cause The cause.
     */
    public ADSException(String message, Throwable cause) {
        super(message, cause);
    }

    /**
     * Constructor to instance an ADSException object with the cause.
     * @param cause The cause.
     */
    public ADSException(Throwable cause) {
        super(cause);
    }
}
