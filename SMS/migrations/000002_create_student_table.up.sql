CREATE TYPE "departmentType" AS ENUM (
    'CSE',
    'IT',
    'ECE',
    'EEE',
    'MECH',
    'CIVIL',
    'AIDS',
    'AIML'
);

CREATE TABLE IF NOT EXISTS student(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    department "departmentType" NOT NULL,
    semester VARCHAR(255) NOT NULL,
    age INTEGER NOT NULL,
    "createdAtUTC" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAtUTC" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);