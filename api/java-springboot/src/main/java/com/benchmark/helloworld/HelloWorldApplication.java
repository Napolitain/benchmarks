package com.benchmark.helloworld;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@SpringBootApplication
@RestController
public class HelloWorldApplication {

    public static void main(String[] args) {
        System.setProperty("server.tomcat.threads.max", "1");
        System.setProperty("server.tomcat.threads.min-spare", "1");
        SpringApplication.run(HelloWorldApplication.class, args);
        System.out.println("Spring Boot server listening on :8080 (single-threaded)");
    }

    @GetMapping("/")
    public Map<String, String> hello() {
        return Map.of("message", "Hello, World!");
    }
}
