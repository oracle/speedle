package oracle.speedle.ads.client;

public class Principal {
    public static enum PrincipalType {
        USER("user"),
        GROUP("group");

        final private String typeName;
        private PrincipalType(String typeName) {
            this.typeName = typeName;
        }

        public String getTypeName() {
            return typeName;
        }

        @Override
        public String toString() {
            return typeName;
        }
    }

    private PrincipalType type = null;
    private String name = null;

    public PrincipalType getType() {
        return type;
    }

    public void setType(PrincipalType type) {
        this.type = type;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    @Override
    public String toString() {
        return "Type: " + type + " Name: " + name;
    }
}
