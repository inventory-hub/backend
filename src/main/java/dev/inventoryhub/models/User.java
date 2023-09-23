package dev.inventoryhub.models;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

import java.sql.Timestamp;

@Data
@Entity
@Table(name = "users")
public class User {
    @Id
    private long user_id;
    private String user_username;
    private String user_email;
    private String user_password;
    private Timestamp user_creation_date;
    private Timestamp user_last_login_date;
}
