package analysis

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/jhillyerd/enmime"
	"github.com/manveru/faker"
)

var (
	fakes *faker.Faker
)

func init() {
	var err error
	fakes, err = faker.New("en")
	if err != nil {
		panic(err)
	}
	// seed with a constant value so we can generate the same fake data each time
	gofakeit.Seed(0)
}

// WritePartsToDisk stores all emails on disk in mime format
func WritePartsToDisk(parts []*enmime.Part, outdir string) error {
	if err := os.Mkdir(outdir, os.ModeDir); !os.IsExist(err) {
		return err
	}
	for i, part := range parts {
		buf := bytes.NewBuffer(nil)
		if err := part.Encode(buf); err != nil {
			return err
		}
		ioutil.WriteFile(
			fmt.Sprintf("%s/email-%v.eml", outdir, i),
			buf.Bytes(),
			os.FileMode(0642),
		)
	}
	return nil
}

// GenerateMessages is used to generate fake email messages
func GenerateMessages(count int) ([]*enmime.Part, error) {
	var emails = make([]*enmime.Part, count)
	addresses := GenerateFakeEmails(count)
	if len(addresses) == 0 {
		return nil, errors.New("failed to generate addresses")
	}
	for i := 0; i < count; i++ {
		part, err := GenerateMessage(fakes, GenOpts{To: addresses[i]}).Build()
		if err != nil {
			return nil, err
		}
		emails[i] = part
	}
	return emails, nil
}

type GenOpts struct {
	To             string
	ParagraphCount int
	Signature      string
}

// GenerateMessage uses faker to create a random message struct
func GenerateMessage(fake *faker.Faker, opts GenOpts) enmime.MailBuilder {
	var to = opts.To
	from := fake.Email()
	company := fake.CompanyName()
	// Plain text
	cosig := fmt.Sprintf("%s <%s>, %s\r\n%s, \"%s\"",
		fake.Name(),
		from,
		fake.JobTitle(),
		company,
		fake.CompanyCatchPhrase())
	paragraphs := fake.Paragraphs(4, true)
	textp := append(make([]string, 0), paragraphs...)
	textp = append(textp, cosig)
	/* TODO(bonedaddy): generate fake signature
	if *signature != "" {
		textp = append(textp, "--\r\n"+*signature)
	}
	*/
	// HTML
	cosig = fmt.Sprintf("%s &lt;<a href=\"mailto:%s\">%s</a>&gt;, %s<br>\r\n<b>%s</b>, <em>%s</em>",
		fake.Name(),
		from,
		from,
		fake.JobTitle(),
		company,
		fake.CompanyCatchPhrase())
	htmlp := append(make([]string, 0), paragraphs...)
	htmlp = append(htmlp, cosig)
	/* TODO(bonedaddy): generate fake signature
	if *signature != "" {
		htmlp = append(htmlp, "<small>"+*signature+"</small>")
	}
	*/
	return enmime.Builder().
		From("", from).
		To("", to).
		Subject(strings.Title(fake.CompanyBs()) + " with " + company).
		Text([]byte(strings.Join(textp, "\r\n\r\n"))).
		HTML([]byte("<p>" + strings.Join(htmlp, "</p>\r\n<p>") + "</p>"))
}

// GenerateFakeEmails is used to generate fake emai laddresses
func GenerateFakeEmails(count int) []string {
	var addresses = make([]string, count)
	for i := 0; i < count; i++ {
		addresses[i] = gofakeit.Email()
	}
	return addresses
}
