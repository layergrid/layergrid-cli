FROM gcr.io/distroless/static-debian12:nonroot

COPY layergrid /layergrid

ENTRYPOINT ["/layergrid"]
CMD ["scan", "/scan"]
