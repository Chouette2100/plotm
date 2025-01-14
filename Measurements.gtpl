<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
<form>
<input type="submit" value="Display previous" formaction="Measurements?fnc=P" formmethod="POST" style="background-color: khaki">
<input type="submit" value="Redisplay" formaction="Measurements?fnc=R" formmethod="POST" style="background-color: khaki">
<input type="submit" value="Display next" formaction="Measurements?fnc=N" formmethod="POST" style="background-color: khaki">
<input style="margin-left: 1em;background-color: pink;" type="submit" value="Display latest" formaction="Measurements?fnc=L" formmethod="POST">
<input style="margin-left: 3em; background-color: aquamarine;" type="submit" value="Restore the settings from" formaction="Measurements?fnc=RS" formmethod="POST" >
<select name="yml">
    {{ range $i, $y := .Ymlfiles }}
    <option value="{{ $y }}">{{ $y }}</option>
    {{ end }}
</select>
<input style="margin-left: 3em; background-color: aquamarine;" type="submit" value="Save the settings in" formaction="Measurements?fnc=SV" formmethod="POST">
<input style="background-color: gainsboro;" type="text" value="{{ .Nextyml }}" readonly style="width: 100px;" />
<input style="margin-left: 1px;" type="checkbox" name="current" value="yes" />
    <label for="current">Set the end time to the time of execution</label>

<br>
<label for="pet-select">Choose the display period:</label>
<select name="period" id="period">
  <option value="">--Please choose an option--</option>
  <option {{ if eq .Period "8 days" }} selected {{ end }} value="8 days">8 days</option>
  <option {{ if eq .Period "4 days" }} selected {{ end }} value="4 days">4 days</option>
  <option {{ if eq .Period "2 days" }} selected {{ end }} value="2 days">2 days</option>
  <option {{ if eq .Period "1 day" }} selected {{ end }} value="1 day">1 day</option>
  <option {{ if eq .Period "12 hours" }} selected {{ end }} value="12 hours">12 hours</option>
  <option {{ if eq .Period "6 hours" }} selected {{ end }} value="6 hours">6 hours</option>
  <option {{ if eq .Period "3 hours" }} selected {{ end }} value="3 hours">3 hours</option>
  <option {{ if eq .Period "2 hours" }} selected {{ end }} value="2 hours">2 hours</option>
  <option {{ if eq .Period "1 hour" }} selected {{ end }} value="1 hour">1 hour</option>
  <option {{ if eq .Period "30 minutes" }} selected {{ end }} value="30 minutes">30 minutes</option>
  <option {{ if eq .Period "15 minutes" }} selected {{ end }} value="20 minutes">5 minutes</option>
  <option {{ if eq .Period "10 minutes" }} selected {{ end }} value="10 minutes">10 minutes</option>
  <option {{ if eq .Period "4 minutes" }} selected {{ end }} value="4 minutes">4 minutes</option>
</select>
{{/*
<input id="method" name="method" type="hidden" value="{{ .Method }}">
*/}}

{{/*<input id="nterm" name="nterm" type="hidden" value="{{ .Nterm }}">*/}}
<input id="uetime" name="uetime" type="hidden" value="{{ .Uetime }}">
<br>
<table>
<tr>
<td>
<img src="{{.Filename}}" alt="" height="100%">
</td>
<td>
{{ $m := .Method }}
{{ range $i, $v := .Item }}
{{ .Name }}
<input type="number" name="rng_{{$i}}_min" value="{{.Vmin}}" style="width: 60px;"
    {{ if eq .Name "CO2" }} step=500 {{ end }}
    {{ if eq .Name "Humidity" }} step=5 {{ end }}
    style="text-align: right; margin-left: 2em;" /> -- 
<input type="number" name="rng_{{$i}}_max" value="{{.Vmax}}" style="width: 60px;"
    {{ if eq .Name "CO2" }} step=500 {{ end }}
    {{ if eq .Name "Humidity" }} step=5 {{ end }}
    style="text-align: right;" />
<br>
{{ if eq .Name "Humidity" }}
<input type="radio" name="method" value="R" {{ if eq $m "R" }}checked{{ end }}>Relative humidity<br>
<input type="radio" name="method" value="V" {{ if eq $m "V" }}checked{{ end }}>Absolute humidity<br>
{{ end }}
{{ range $j, $w := .Udev }}
    <input type="checkbox" name="dev{{$i}}_{{$j}}" value="checked" {{.Status}} >{{ .Name }}<br>
{{ end }}
<br><br>
{{ end }}
{{/*
Temperature<br>
{{ range .Device }}
    <input type="checkbox" name="dev{{.Device}}" value="checked" {{.Status}} >{{ .Name }}<br>
{{ end }}
<br>
<br>Humidity<br>
<input type="radio" name="method" value="V" {{ if eq .Method "V" }}checked{{ end }}>Absolute humidity<br>
<input type="radio" name="method" value="R" {{ if eq .Method "R" }}checked{{ end }}>Relative humidity<br>
{{ range .Device }}
    <input type="checkbox" name="dev{{.Device}}" value="checked" {{.Status}} >{{ .Name }}<br>
{{ end }}
<br>CO2<br>
{{ range .Device }}
    <input type="checkbox" name="dev{{.Device}}" value="checked" {{.Status}} >{{ .Name }}<br>
{{ end }}
*/}}
<br>
</td>
</tr>
</table>
The source of this program can be found <a href="https://github.com/Chouette2100/plotm">here</a>
</form>
</body>
</html>
