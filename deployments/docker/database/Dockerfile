FROM postgres:15-alpine

# Set non-root user (postgres user created by default in this image)
USER postgres

# Initialize DB with pgcrypto for encrypted columns/passwords
# The initdb.d scripts are run automatically during database initialization
COPY ./init-scripts/ /docker-entrypoint-initdb.d/

# Restrict permissions on the entrypoint scripts
USER root
RUN chmod -R 750 /docker-entrypoint-initdb.d/ && \
    chown -R postgres:postgres /docker-entrypoint-initdb.d/
USER postgres

# Expose standard port
EXPOSE 5432
