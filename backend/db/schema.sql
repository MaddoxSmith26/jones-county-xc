-- Jones County XC Database Schema

CREATE DATABASE IF NOT EXISTS jones_county_xc;
USE jones_county_xc;

-- Athletes table
CREATE TABLE IF NOT EXISTS athletes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    grade INT NOT NULL,
    personal_record VARCHAR(10),
    events VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Meets table
CREATE TABLE IF NOT EXISTS meets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    date DATE NOT NULL,
    location VARCHAR(200),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Results table (links athletes to meets)
CREATE TABLE IF NOT EXISTS results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    athlete_id INT NOT NULL,
    meet_id INT NOT NULL,
    time VARCHAR(10) NOT NULL,
    place INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (athlete_id) REFERENCES athletes(id) ON DELETE CASCADE,
    FOREIGN KEY (meet_id) REFERENCES meets(id) ON DELETE CASCADE
);
