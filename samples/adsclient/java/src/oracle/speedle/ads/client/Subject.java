package oracle.speedle.ads.client;

import java.util.List;

public class Subject {
    private List<Principal> principals = null;

    public Subject() {}

    public void setPrincipals(List<Principal> principal) {
        this.principals = principal;
    }
    public List<Principal> getPrincipals() {
        return principals;
    }

    @Override
    public String toString() {
        return "" + principals;
    }
}
