FROM golang:alpine AS build
ADD api.go .
RUN go build -v -o /api api.go

FROM build as repo
RUN apk --no-cache add git
RUN git clone https://github.com/koenbuyens/securityheaders /src

FROM python:3-alpine
WORKDIR /app

COPY --from=repo /src/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY --from=repo /src/ .
RUN python securityheaders.py --formatter json google.com > /dev/null

COPY --from=build /api .
EXPOSE 8080
CMD /app/api
