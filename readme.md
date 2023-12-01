
`nerdctl build . -t gotab:dev`

`nerdctl run -p 3000:8080 -v /tmp/gotab:/gotab-disk --rm gotab:dev --disk /gotab-disk --verbosity 4`