package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"crypto/tls"
	"crypto/x509"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/ocsp"
)


const	(
	TLSMaxConn	int	= 1000
)


var	(
	TLS_client_session_cache	tls.ClientSessionCache	=	tls.NewLRUClientSessionCache(TLSMaxConn)
)

type	(

	TLSLevel	int

	TLSClientConfig struct {
		Strength		TLSLevel
		// ServerName for SNI negociation
		SNI			string

		// this is the list of BASE64 SHA256 fingerprints of the already known PublicKey in the certificate
		// true means known as valid.
		// false means known as invalid.
		// unset means unknown.
		PKPs			map[string]bool

		handshake_done		bool
		rootCAs			*x509.CertPool
		clientCAs		*x509.CertPool
	}


	TLSState struct {
		OCSPExist	bool
		OCSPValid	bool
		OCSPUnknown	bool
		OCSPChecked	bool

		SNIExist	bool
		SNIValid	bool

		PKPExist	bool
		PKPValid	int
		PKPInvalid	int
		PKPCerts	int
	}
)

const	(
	TLS_SECURE			TLSLevel = iota
	TLS_I_ACCEPT_SOME_PRIVACY_RISK_AND_I_TAKE_THE_RISK_OF_A_CLASS_ACTION
	TLS_I_DONT_WANT_ANY_PRIVACY_AND_I_LOVE_CLASS_ACTION_IN_COURT
)




func (tcc *TLSClientConfig) GetTLSConfig() *tls.Config {
	switch	tcc.Strength {
	case	TLS_I_ACCEPT_SOME_PRIVACY_RISK_AND_I_TAKE_THE_RISK_OF_A_CLASS_ACTION:
		return	&tls.Config {
			ServerName:		tcc.SNI,
			MinVersion:		tls.VersionTLS10,
			MaxVersion:		tls.VersionTLS12,
			ClientSessionCache:	TLS_client_session_cache,
			CurvePreferences:	[]tls.CurveID {
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256,
			},
			CipherSuites:		[]uint16 {
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

				// Sensible to Lucky13 attack
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			},
		}

	case	TLS_I_DONT_WANT_ANY_PRIVACY_AND_I_LOVE_CLASS_ACTION_IN_COURT:
		return	&tls.Config {
			ServerName:		tcc.SNI,
			MinVersion:		tls.VersionSSL30,
			MaxVersion:		tls.VersionTLS12,
			ClientSessionCache:	TLS_client_session_cache,
			CurvePreferences:	[]tls.CurveID {
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256,
			},
			CipherSuites:		[]uint16 {
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,

				tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,

				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_RSA_WITH_RC4_128_SHA,
			},
		}

	default:
		return	&tls.Config {
			ServerName:		tcc.SNI,
			MinVersion:		tls.VersionTLS12,
			MaxVersion:		tls.VersionTLS12,
			ClientSessionCache:	TLS_client_session_cache,
			CurvePreferences:	[]tls.CurveID {
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256,
			},
			CipherSuites:		[]uint16 {
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
	}
}









func (tcc *TLSClientConfig) Verify(conn *tls.Conn) (*TLSState, error) {
	var ocsprep	*ocsp.Response
	var der		[]byte
	var err		error

	res	:= new(TLSState)
	cstate	:= conn.ConnectionState()

	res.SNIExist	= (tcc.SNI != "")
	res.PKPExist	= (tcc.PKPs != nil && len(tcc.PKPs) > 0)

	if cstate.OCSPResponse != nil {
		ocsprep,err	= ocsp.ParseResponse( cstate.OCSPResponse, nil )
		if err != nil {
			return	nil, err
		}
		res.OCSPExist	= true
		res.OCSPValid	= (ocsprep.Status == ocsp.Good)
		res.OCSPUnknown	= (ocsprep.Status == ocsp.Unknown)
	}

	for _, peercert := range cstate.PeerCertificates {
		der,err	= x509.MarshalPKIXPublicKey(peercert.PublicKey)
		if err != nil {
			return	nil, err
		}

		if res.SNIExist && !res.SNIValid && peercert.VerifyHostname( tcc.SNI ) == nil {
			res.SNIValid = true
		}

		if res.OCSPValid && !res.OCSPChecked && ocsprep.CheckSignatureFrom(peercert) == nil {
			res.OCSPChecked	= true
		}

		rawhash	:= sha256.Sum256(der)
		hash	:= base64.StdEncoding.EncodeToString( rawhash[:] )

		if res.PKPExist {
			res.PKPCerts++
			valid, ok := tcc.PKPs[hash]
			switch	{
			case	ok && valid:
				res.PKPValid++
			case	ok && !valid:
				res.PKPInvalid++
			}
		}
	}

	return	res, nil
}


func (t TLSState) SNI() bool {
	return	!t.SNIExist ||
		t.SNIValid
}

func (t TLSState) OCSPstrict() bool {
	return	!t.OCSPExist ||
		(t.OCSPValid && t.OCSPChecked)
}

func (t TLSState) PKPstrict() bool {
	return	!t.PKPExist ||
		(t.PKPInvalid==0)
}

func (t TLSState) IsTrustableStrict() bool {
	return t.SNI() && t.OCSPstrict() && t.PKPstrict()
}
