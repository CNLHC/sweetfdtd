package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

const normal_template = `
<addcolumn><column><title>Max time remaining</title><type>text</type><tagname>time</tagname></column><column><title>Autoshutoff level</title><type>text</type><tagname>shutoff</tagname></column><column><title>Progress</title><type>progress</type><tagname>progress</tagname></column></addcolumn>
<propertylist><property>canBePaused</property><property>canQuitAndSave</property></propertylist>
<status>UPDATE_MODES</status>
<status>MESH_START</status>
<status>BUILD_SOURCES_START</status>
<status>BUILD_SOURCES_FINISH</status>
0.0491584,20 mins, 1 secs,1,
2.08104,3 mins, 44 secs,0.970258,
<update><progress>3.65411</progress><time>2 mins, 39 secs</time><shutoff>1.01144</shutoff></update>
<update><progress>5.75153</progress><time>1 min, 43 secs</time><shutoff>0.994792</shutoff></update>
<update><progress>9.94638</progress><time>1 min, 28 secs</time><shutoff>0.976302</shutoff></update>
<update><progress>12.0438</progress><time>1 min, 21 secs</time><shutoff>0.975366</shutoff></update>
<update><progress>14.1412</progress><time>1 min, 11 secs</time><shutoff>1.00209</shutoff></update>
<update><progress>18.3361</progress><time>1 min, 3 secs</time><shutoff>1.00776</shutoff></update>
<update><progress>20.4335</progress><time>57 secs</time><shutoff>0.989772</shutoff></update>
<update><progress>24.6283</progress><time>52 secs</time><shutoff>0.98139</shutoff></update>
<update><progress>26.7258</progress><time>48 secs</time><shutoff>0.994363</shutoff></update>
<update><progress>30.9206</progress><time>43 secs</time><shutoff>0.998544</shutoff></update>
<update><progress>35.1155</progress><time>47 secs</time><shutoff>0.991905</shutoff></update>
<update><progress>37.2129</progress><time>47 secs</time><shutoff>1.00439</shutoff></update>
<update><progress>38.2616</progress><time>46 secs</time><shutoff>0.989403</shutoff></update>
<update><progress>40.359</progress><time>55 secs</time><shutoff>1.0081</shutoff></update>
<update><progress>41.4077</progress><time>53 secs</time><shutoff>1.00049</shutoff></update>
<update><progress>43.5052</progress><time>56 secs</time><shutoff>0.991988</shutoff></update>
<update><progress>44.5539</progress><time>55 secs</time><shutoff>1.01107</shutoff></update>
<update><progress>46.6513</progress><time>1 min, 15 secs</time><shutoff>0.998264</shutoff></update>
<update><progress>47.7</progress><time>1 min, 14 secs</time><shutoff>0.994334</shutoff></update>
<update><progress>49.2731</progress><time>1 min, 12 secs</time><shutoff>1.0068</shutoff></update>
<update><progress>50.3218</progress><time>1 min, 9 secs</time><shutoff>0.993763</shutoff></update>
<update><progress>52.4192</progress><time>1 min, 5 secs</time><shutoff>0.996724</shutoff></update>
<update><progress>54.5166</progress><time>1 min, 0 secs</time><shutoff>1.00381</shutoff></update>
<update><progress>56.6141</progress><time>57 secs</time><shutoff>0.982359</shutoff></update>
<update><progress>58.7115</progress><time>54 secs</time><shutoff>0.962238</shutoff></update>
<update><progress>59.7602</progress><time>53 secs</time><shutoff>0.968263</shutoff></update>
<update><progress>60.8089</progress><time>51 secs</time><shutoff>0.954698</shutoff></update>
<update><progress>62.9063</progress><time>47 secs</time><shutoff>0.990783</shutoff></update>
<update><progress>65.0038</progress><time>44 secs</time><shutoff>0.999365</shutoff></update>
<update><progress>67.1012</progress><time>40 secs</time><shutoff>0.976683</shutoff></update>
<update><progress>71.296</progress><time>33 secs</time><shutoff>0.987137</shutoff></update>
<update><progress>75.4909</progress><time>27 secs</time><shutoff>0.983458</shutoff></update>
<update><progress>83.8806</progress><time>16 secs</time><shutoff>0.991594</shutoff></update>
<update><progress>92.2703</progress><time>7 secs</time><shutoff>0.9873</shutoff></update>
<update><progress>100</progress><time>0 secs</time><shutoff>0.980719</shutoff></update>
<status>WRITE_START</status>
<status>WRITE_FINISH</status>
<simComplete/>
<complete/>
`
const error_template = `
<addcolumn><column><title>Max time remaining</title><type>text</type><tagname>time</tagname></column><column><title>Autoshutoff level</title><type>text</type><tagname>shutoff</tagname></column><column><title>Progress</title><type>progress</type><tagname>progress</tagname></column></addcolumn>
<propertylist><property>canBePaused</property><property>canQuitAndSave</property></propertylist>
c06b02n08(process 0): Your license settings don't appear to be configured.
Please open the Launcher or the Configure License program to reconfigure your license settings.
c06b02n08(process 0): License error: The license settings have not been configured<p>The flexNet error code is: -4, which corresponds to the error:</p></p>Licensed number of users already reached.
Feature:       FDTD_Solutions_engine
License path:  27011@11.3.11.1:27011@localhost:
FlexNet Licensing error:-4,132.  System Error: 2 "No such file or directory"</p><p>Please see <a href=https://kb.lumerical.com/redirect/fwd3.html>Troubleshooting guide</a> for help resolving this issue.</p>
c06b02n08(process 0): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 10): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 11): Error: there was a failure with the license. Process number: 0 had this error
<complete/>
c06b02n08(process 16): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 25): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 17): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 27): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 26): Error: there was a failure with the license. Process number: 0 had this errorc06b02n08(process 9): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 15): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 21): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 23): Error: there was a failure with the license. Process number: 0 had this error

c06b02n08(process 8): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 12): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 19): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 1): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 4): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 3): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 14): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 5): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 7): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 13): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 22): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 18): Error: there was a failure with the license. Process number: 0 had this errorc06b02n08(process 2): Error: there was a failure with the license. Process number: 0 had this error

c06b02n08(process 24): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 20): Error: there was a failure with the license. Process number: 0 had this error
c06b02n08(process 6): Error: there was a failure with the license. Process number: 0 had this error
`

func main() {
	var template string
	rand.Seed(time.Now().Unix())

	if rand.Intn(100) >= 50 {
		template = normal_template
	} else {
		template = error_template
	}
	lines := strings.Split(template, "\n")
	writer := bufio.NewWriter(os.Stdout)

	for _, line := range lines {
		writer.WriteString(line + "\n")
		writer.Flush()
		time.Sleep(time.Millisecond * time.Duration(0+rand.Intn(100)))
	}

}
