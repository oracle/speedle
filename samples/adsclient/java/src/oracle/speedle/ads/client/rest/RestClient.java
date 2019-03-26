package oracle.speedle.ads.client.rest;

import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStreamWriter;
import java.lang.reflect.Type;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.URL;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.HashMap;
import java.util.Map;
import java.util.logging.Level;
import java.util.logging.Logger;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.JsonElement;
import com.google.gson.JsonPrimitive;
import com.google.gson.JsonSerializationContext;
import com.google.gson.JsonSerializer;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;

import oracle.speedle.ads.client.ADSException;
import oracle.speedle.ads.client.Client;
import oracle.speedle.ads.client.RequestContext;
import oracle.speedle.ads.client.Principal.PrincipalType;

public class RestClient implements Client {
    private static final String CLASS_NAME = RestClient.class.getName();
    private static final Logger LOGGER = Logger.getLogger(CLASS_NAME);

    private static Gson gson;
    
    static {
        gson = new GsonBuilder().registerTypeAdapter(PrincipalType.class, new PrincipalTypeSerializer())
                .registerTypeAdapter(ZonedDateTime.class, new ZonedDateTimeSerializer())
                .create();
    }

    private static class PrincipalTypeSerializer implements JsonSerializer<PrincipalType> {
        public JsonElement serialize(PrincipalType pt, Type typeOfSrc, JsonSerializationContext context) {
            return new JsonPrimitive(pt.toString());
        }
    }

    private static class ZonedDateTimeSerializer implements JsonSerializer<ZonedDateTime> {
        private static DateTimeFormatter dateTimeFormatter = DateTimeFormatter.ofPattern("yyyy-MM-dd'T'HH:mm:ss.SSS");

        public JsonElement serialize(ZonedDateTime dt, Type typeOfSrc, JsonSerializationContext context) {
            return new JsonPrimitive(dt.format(dateTimeFormatter) + dt.getOffset().getId());
        }
    }

    private final URL isAllowedEndpoint;

    public RestClient(String endpoint) throws ADSException {
        if (endpoint == null) {
            throw new IllegalArgumentException("host is null.");
        }

        try {
            isAllowedEndpoint = new URL(String.format("http://%s/authz-check/v1/is-allowed", endpoint));
        } catch (MalformedURLException e) {
            throw new ADSException(e);
        }
    }

    @Override
    public boolean isAllowed(RequestContext context) throws ADSException {
        final String methodName = "isAllowed";
        if (LOGGER.isLoggable(Level.FINER)) {
            LOGGER.entering(CLASS_NAME, methodName, context);
        }

        Map<String, String> headers = new HashMap<String, String>();
        headers.put("Content-Type", "application/json");

        HttpURLConnection connection = post(isAllowedEndpoint, headers, context);

        JsonReader reader = null;
        try {
            int responseCode = connection.getResponseCode();
            if (responseCode != 200) {
                throw new ADSException("Server returns unexpected status code " + responseCode + ".");
            }

            reader = gson.newJsonReader(new InputStreamReader(connection.getInputStream()));
            IsAllowedResponse outputObj = gson.fromJson(reader, IsAllowedResponse.class);

            if (LOGGER.isLoggable(Level.FINER)) {
                LOGGER.exiting(CLASS_NAME, methodName);
            }
            return outputObj.isAllowed();
        } catch (IOException e) {
            throw new ADSException(e);
        } finally {
            if (reader != null) {
                try { reader.close(); } catch (IOException e) {
                    String msg = "Exception in closing input stream.";
                    LOGGER.logp(Level.WARNING, CLASS_NAME, methodName, msg, e);
                }
            }
        }
    }

    private static HttpURLConnection post(URL url, Map<String, String> headers, RequestContext context) throws ADSException {
        final String methodName = "get";

        if (LOGGER.isLoggable(Level.FINER)) {
            LOGGER.entering(CLASS_NAME, methodName, url);
        }

        HttpURLConnection conn = null;
        try {
            conn = (HttpURLConnection)url.openConnection();
            conn.setUseCaches(false);
            conn.setDoOutput(true);
            conn.setRequestMethod("POST");
            if (headers != null) {
                for (Map.Entry<String, String> header : headers.entrySet()) {
                    conn.setRequestProperty(header.getKey(), header.getValue());
                }
            }
            conn.connect();

            JsonWriter writer = gson.newJsonWriter(new OutputStreamWriter(conn.getOutputStream()));
            gson.toJson(context, RequestContext.class, writer);
            writer.flush();
        } catch (IOException e) {
            throw new ADSException(e);
        } finally {
            if (conn != null) {
                conn.disconnect();
            }
        }

        if (LOGGER.isLoggable(Level.FINER)) {
            LOGGER.exiting(CLASS_NAME, methodName, conn);
        }
        return conn;
    }
}
