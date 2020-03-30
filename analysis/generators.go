package analysis

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"image"
	"image/color"
	"image/png"
	"os"

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

// GenerateMessages is used to generate fake email messages and save to disk
func GenerateMessages(outdir string, count, numEmojis, paragraphCount int) error {
	os.MkdirAll(outdir, os.ModePerm)
	addresses := GenerateFakeEmails(count)
	if len(addresses) == 0 {
		return errors.New("failed to generate addresses")
	}
	for i := 0; i < count; i++ {
		part := GenerateMessage(fakes, GenOpts{To: addresses[i], EmojiCount: numEmojis, ParagraphCount: paragraphCount})
		part.AddAttachment(genImage(), "image/png", fmt.Sprintf("image-%v.png", i))
		email, err := part.Build()
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(nil)
		if err := email.Encode(buf); err != nil {
			return err
		}
		if err := ioutil.WriteFile(
			fmt.Sprintf("%s/email-%v.eml", outdir, i),
			buf.Bytes(),
			os.FileMode(0642),
		); err != nil {
			return err
		}
	}
	return nil
}

// GenOpts is used to control generation of email messages
type GenOpts struct {
	To             string
	ParagraphCount int
	Signature      string
	EmojiCount     int
}

// GenerateMessage uses faker to create a random message struct
func GenerateMessage(fake *faker.Faker, opts GenOpts) enmime.MailBuilder {
	var to = opts.To
	from := fake.Email()
	company := fake.CompanyName()
	// Plain text
	cosig := fmt.Sprintf(
		"%s <%s>, %s\r\n%s, \"%s\"\nheres a bunch of emojies %s\n%s",
		fake.Name(),
		from,
		fake.JobTitle(),
		company,
		fake.CompanyCatchPhrase(),
		emojiSpam(opts.EmojiCount),
		gofakeit.BS(),
	)
	paragraphs := fake.Paragraphs(opts.ParagraphCount, true)
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

func emojiSpam(count int) string {
	var spam string
	for i := 0; i < count; i++ {
		spam = fmt.Sprintf("%s-%s", spam, gofakeit.Emoji())
	}
	return spam
}

// GenerateFakeEmails is used to generate fake emai laddresses
func GenerateFakeEmails(count int) []string {
	var addresses = make([]string, count)
	for i := 0; i < count; i++ {
		addresses[i] = gofakeit.Email()
	}
	return addresses
}

func genImage() []byte {
	buf := bytes.NewBuffer(nil)
	img := image.NewRGBA(image.Rect(0, 0, 1920, 1080))

	// Draw a red dot at (2, 3)
	img.Set(2, 3, color.RGBA{uint8(rand.Int63n(255)), uint8(rand.Int63n(255)), uint8(rand.Int63n(255)), uint8(rand.Int63n(255))})
	if err := png.Encode(buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
