import org.yaml.snakeyaml.Yaml;
import java.io.FileInputStream;
import java.io.InputStream;
import java.util.Map;

public class Rectangle {
    public static void main(String[] args) {
        if (args.length != 1) {
            System.err.println("Usage: java Rectangle <yaml-file>");
            System.exit(1);
        }

        String yamlFile = args[0];
        long startTime = System.nanoTime();

        try (InputStream inputStream = new FileInputStream(yamlFile)) {
            Yaml yaml = new Yaml();
            Map<String, Object> data = yaml.load(inputStream);

            double a = ((Number) data.getOrDefault("a", 0)).doubleValue();
            double b = ((Number) data.getOrDefault("b", 0)).doubleValue();
            double c = ((Number) data.getOrDefault("c", 0)).doubleValue();
            double d = ((Number) data.getOrDefault("d", 0)).doubleValue();

            double width = Math.abs(c - a);
            double height = Math.abs(d - b);
            double area = width * height;

            long endTime = System.nanoTime();
            double elapsedMs = (endTime - startTime) / 1_000_000.0;

            System.out.printf("Rectangle area: %.2f%n", area);
            System.out.printf("Time: %.6f ms%n", elapsedMs);

        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
            System.exit(1);
        }
    }
}
