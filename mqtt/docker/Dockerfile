FROM elixir:1.9.1-alpine AS plugin_builder

# Set environment variables for building the application    

ENV MIX_ENV=prod \
    TEST=1 \
    LANG=C.UTF-88

RUN apk add build-base

RUN apk add --update git && \
    rm -rf /var/cache/apk/*

# Install hex and rebar
RUN mix local.hex --force && \
    mix local.rebar --force

# Create the application build directory
RUN mkdir /app
WORKDIR /app

ADD http://api.geonames.org/timezoneJSON?formatted=true&lat=47.01&lng=10.2&username=demo&style=full /tmp/bustcache

# create child app directories
RUN mkdir -p vernemq_iot_plugin/

COPY ./vernemq_iot_plugin/ vernemq_iot_plugin/
WORKDIR /app/vernemq_iot_plugin/

RUN mix deps.get
RUN mix deps.compile
RUN mix release

FROM vernemq/vernemq:1.10.3-alpine

WORKDIR /home/app

COPY --from=plugin_builder /app/vernemq_iot_plugin/_build .

COPY ./docker/vernemq.conf /etc/vernemq/
