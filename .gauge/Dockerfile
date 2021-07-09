# This image uses the official node base image.
FROM node
 
# Set a custom npm install location so that Gauge, Taiko and dependencies can be 
# installed without root privileges
ENV NPM_CONFIG_PREFIX=/home/gauge/.npm-packages
ENV PATH="${NPM_CONFIG_PREFIX}/bin:${PATH}"

# Set working directory
WORKDIR /gauge
 
# Copy the local working folder
COPY . .

# Create an unprivileged user
RUN groupadd -r gauge && useradd -r -g gauge -G audio,video gauge && \
   mkdir -p /home/gauge/.npm-packages/lib && \
   chown -R gauge:gauge /home/gauge /gauge
 
USER gauge

# Install dependencies and plugins
RUN npm install -g @getgauge/cli \
    && npm install \
    && gauge install \
    && gauge config check_updates false

# Default command on running the image
ENTRYPOINT ["npm", "test"]
