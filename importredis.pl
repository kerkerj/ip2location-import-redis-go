use strict;
use Redis;

#my $filename = 'IP-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ZIPCODE-TIMEZONE-ISP-DOMAIN-NETSPEED-AREACODE-WEATHER-MOBILE-ELEVATION-USAGETYPE.CSV'; # IP2Location DB24
my $filename = './testdata/IPV6-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ISP-DOMAIN-MOBILE-USAGETYPE.CSV'; # IP2Location DB23

my $data = '';
my $dbname = 'DB23';
my $client = Redis->new( server => 'localhost:6379' );
$client->del($dbname);
open IN, "<$filename" or die;
my $counter = 0;
my $line = '';
while (<IN>)
{
    $line = $_;
    $line =~ s/[\r\n]+//; # remove EOL
    if ($line =~ /^"([^"]+)","[^"]+","(.+)"$/)
    {
        $counter++;
        my $ipfrom = $1;
        my $others = $2;
        $others =~ s/","/\|/g;
        my $datastr = '"' . $ipfrom . "|" . $others . '"';
        $datastr =~ s/([\@\%\$])/\\\1/g;
        $data .= ', ' . $ipfrom . ' => ' . $datastr;

        if ($counter % 5000 == 0)
        {
            print "Running line $counter\n";
            &runMe($data);
            $data = '';
        }
    }
}

if ($data ne '')
{
    &runMe($data);
    $data = '';
}
close IN;

sub runMe
{
    my $mydata = shift;
    eval('$client->zadd($dbname' . $mydata . ');');
    warn $@ if $@;
}
