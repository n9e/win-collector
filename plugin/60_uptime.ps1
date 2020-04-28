# get localIP
$ips = foreach($ipv4 in (ipconfig) -like '*IPv4*') { ($ipv4 -split ' : ')[-1]}
if ($ips -is [array]){
	$localIP = $ips[0]
}else{
	$localIP = $ips
}


# get timestamp
$ts = [int64](Get-Date(Get-Date).ToUniversalTime() -UFormat "%s")

# get uptime
$millisec = [Environment]::TickCount
$uptime  = [Timespan]::FromMilliseconds($millisec).TotalSeconds

# get step
$step = ($MyInvocation.MyCommand.Name -split '_')[0]

Write-Host "
[
    {
        `"endpoint`": `"$localIP`",
        `"tags`": `"`",
        `"timestamp`": $ts,
        `"metric`": `"sys.uptime.duration`",
        `"value`": $uptime,
        `"step`": $step
    }
]"