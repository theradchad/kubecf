package main

// credhub_info.go contains functions to determine the credhub information
// (addresses and port) required for our security groups.

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"
)

// credhubLinkData is the quarks entanglement data structure for the credhub link
type credhubLinkData struct {
	CredHub struct {
		InternalURL string `json:"internal_url"`
		Port        int    `json:"port"`
	} `json:"credhub"`
}

// resolveCredHubAddrsGivenLink returns the IP addresses of the credhub service.
func resolveCredHubAddrsGivenLink(ctx context.Context, link credhubLinkData) ([]string, error) {
	credhubURL, err := url.Parse(link.CredHub.InternalURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse credhub link URL %s: %w",
			link.CredHub.InternalURL, err)
	}
	var addrs []string
	fmt.Fprintf(os.Stderr, "Looking up credhub service %s", credhubURL.Hostname())
	for {
		addrs, err = net.DefaultResolver.LookupHost(ctx, credhubURL.Hostname())
		var dnsError *net.DNSError
		if !errors.As(err, &dnsError) {
			break
		}
		if !(dnsError.Temporary() || dnsError.IsNotFound) {
			// Unexpected DNS error; report and die
			return nil, fmt.Errorf("error looking up host %s: %w", credhubURL.Hostname(), err)
		}
		// If CredHub has not finished starting up, DNS resolution will fail
		// (even though the credhub service address is fixed); wait a bit and
		// try again until we succeed (to not blow through our retry quota and
		// end up in CrashLoopBackoff).
		fmt.Fprintf(os.Stderr, ".")
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("could not lookup credhub host %s: %w",
			credhubURL.Hostname(), err)
	}
	fmt.Fprintf(os.Stderr, "\nFound credhub service address: %v\n", addrs)

	return addrs, nil
}

// resolveCredHubInfo returns the IP addresses of the CredHub service and the port.
func resolveCredHubInfo(ctx context.Context) ([]string, int, error) {
	var link credhubLinkData
	err := resolveLink(ctx, "credhub", &link)
	if err != nil {
		return nil, 0, err
	}
	addrs, err := resolveCredHubAddrsGivenLink(ctx, link)
	if err != nil {
		return nil, 0, fmt.Errorf("could not resolve credhub address: %w", err)
	}
	return addrs, link.CredHub.Port, nil
}
