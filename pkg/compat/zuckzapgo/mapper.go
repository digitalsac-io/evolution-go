package zuckzapgo

import (
	"fmt"
	"strings"

	"github.com/EvolutionAPI/evolution-go/pkg/utils"
	send_service "github.com/EvolutionAPI/evolution-go/pkg/sendMessage/service"
)

// MapZuckTextToSendText — MessageText (Body, Phone, ContextInfo, Presence).
func MapZuckTextToSendText(d *ZuckTextDTO) *send_service.TextStruct {
	num := d.ZuckBase.EffectiveRecipient()
	body := FirstNonEmpty(d.Text, d.Message, d.Body, d.BodyPascal)
	out := &send_service.TextStruct{Text: body}
	ApplyZuckBaseToTextStruct(out, &d.ZuckBase, num)
	return out
}

func pickMediaURL(d *ZuckMediaDTO, mediaType string) string {
	switch strings.ToLower(mediaType) {
	case "image":
		return FirstNonEmpty(d.URL, d.Image, d.ImageCap)
	case "audio":
		return FirstNonEmpty(d.URL, d.Audio, d.AudioCap)
	case "video", "ptv":
		return FirstNonEmpty(d.URL, d.Video, d.VideoCap)
	case "document":
		return FirstNonEmpty(d.URL, d.Document, d.DocCap)
	default:
		return FirstNonEmpty(d.URL, d.Image, d.ImageCap, d.Audio, d.AudioCap, d.Video, d.VideoCap, d.Document, d.DocCap)
	}
}

func pickMediaCaption(d *ZuckMediaDTO) string {
	return FirstNonEmpty(d.Caption, d.CaptionCap)
}

func pickFilename(d *ZuckMediaDTO) string {
	return FirstNonEmpty(d.Filename, d.FileNameCap)
}

// MapZuckMediaDTO maps Zuck media payloads to MediaStruct.
func MapZuckMediaDTO(d *ZuckMediaDTO, evolutionType string) *send_service.MediaStruct {
	num := d.ZuckBase.EffectiveRecipient()
	mt := evolutionType
	if strings.EqualFold(mt, "video") && d.IsPtv {
		mt = "ptv"
	}
	url := pickMediaURL(d, mt)
	if strings.TrimSpace(d.Base64) != "" && url == "" {
		url = d.Base64
	}
	out := &send_service.MediaStruct{
		Url:     url,
		Type:    mt,
		Caption: pickMediaCaption(d),
	}
	if fn := pickFilename(d); fn != "" {
		out.Filename = fn
	}
	applyZuckBaseToMedia(out, &d.ZuckBase, num)
	return out
}

// MapZuckStickerToSendSticker — MessageSticker (Sticker / Phone / Id).
func MapZuckStickerToSendSticker(d *ZuckStickerDTO) *send_service.StickerStruct {
	num := d.ZuckBase.EffectiveRecipient()
	st := FirstNonEmpty(d.Sticker, d.StickerCap, d.URL)
	out := &send_service.StickerStruct{Sticker: st}
	applyZuckBaseToSticker(out, &d.ZuckBase, num)
	return out
}

// MapZuckLocationToSendLocation — MessageLocation; se address vazio, usa name (OpenAPI não exige address; Evolution sim).
func MapZuckLocationToSendLocation(d *ZuckLocationDTO) *send_service.LocationStruct {
	num := d.ZuckBase.EffectiveRecipient()
	lat := d.Latitude
	if lat == 0 && d.LatitudeCap != 0 {
		lat = d.LatitudeCap
	}
	lng := d.Longitude
	if lng == 0 && d.LongitudeCap != 0 {
		lng = d.LongitudeCap
	}
	name := FirstNonEmpty(d.Name, d.NameCap)
	addr := d.Address
	if strings.TrimSpace(addr) == "" {
		addr = name
	}
	out := &send_service.LocationStruct{
		Name:      name,
		Latitude:  lat,
		Longitude: lng,
		Address:   addr,
	}
	applyZuckBaseToLocation(out, &d.ZuckBase, num)
	return out
}

// MapZuckContactToSendContact — MessageContact (Phone, Name, Vcard string) + formato Evolution.
func MapZuckContactToSendContact(d *ZuckContactDTO) *send_service.ContactStruct {
	num := d.ZuckBase.EffectiveRecipient()
	vc := utils.VCardStruct{}

	if strings.TrimSpace(d.VcardRaw) != "" {
		fn, ph, org := ParseSimpleVCard(d.VcardRaw)
		vc.FullName, vc.Phone, vc.Organization = fn, ph, org
	}
	if d.Vcard != nil {
		if d.Vcard.FullName != "" {
			vc.FullName = d.Vcard.FullName
		}
		if d.Vcard.Organization != "" {
			vc.Organization = d.Vcard.Organization
		}
		if d.Vcard.Phone != "" {
			vc.Phone = d.Vcard.Phone
		}
	}
	if d.FullName != "" {
		vc.FullName = d.FullName
	}
	if d.NameCap != "" {
		vc.FullName = FirstNonEmpty(vc.FullName, d.NameCap)
	}
	if d.Organization != "" {
		vc.Organization = d.Organization
	}
	if d.ContactPhone != "" {
		vc.Phone = d.ContactPhone
	}

	out := &send_service.ContactStruct{Vcard: vc}
	applyZuckBaseToContact(out, &d.ZuckBase, num)
	return out
}

// MapZuckListToSendList — MessageList (Text, Title opcional, Sections com Rows / RowId).
func MapZuckListToSendList(d *ZuckListDTO) *send_service.ListStruct {
	num := d.ZuckBase.EffectiveRecipient()
	title := FirstNonEmpty(d.Title, d.TitleCap)
	if strings.TrimSpace(title) == "" {
		title = " "
	}
	desc := FirstNonEmpty(d.Description, d.Text, d.TextCap, d.Body)
	btn := FirstNonEmpty(d.ButtonText, d.ButtonTextCap)
	footer := FirstNonEmpty(d.FooterText, d.Footer, d.FooterCap)

	secs := d.effectiveSections()
	sections := make([]send_service.Section, 0, len(secs))
	for _, sec := range secs {
		rowsIn := effectiveSectionRows(sec)
		rows := make([]send_service.Row, 0, len(rowsIn))
		for _, r := range rowsIn {
			rdesc := effectiveRowDescription(r)
			if strings.TrimSpace(rdesc) == "" {
				rdesc = " "
			}
			rows = append(rows, send_service.Row{
				Title:       effectiveRowTitle(r),
				Description: rdesc,
				RowId:       effectiveRowID(r),
			})
		}
		sections = append(sections, send_service.Section{
			Title: effectiveSectionTitle(sec),
			Rows:  rows,
		})
	}

	out := &send_service.ListStruct{
		Title:       title,
		Description: desc,
		ButtonText:   btn,
		FooterText:   footer,
		Sections:     sections,
	}
	applyZuckBaseToList(out, &d.ZuckBase, num)
	return out
}

func normalizeZuckButtonType(t string) string {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "quick_reply", "reply":
		return "reply"
	case "cta_url", "url":
		return "url"
	case "cta_call", "call":
		return "call"
	case "cta_copy", "copy":
		return "copy"
	case "pix", "payment_info", "pix_payment":
		return "pix"
	case "review_and_pay":
		return "review_and_pay"
	default:
		return strings.ToLower(strings.TrimSpace(t))
	}
}

func displayTextForZuckButton(b ZuckInteractiveButtonDTO) string {
	if b.ButtonText.DisplayText != "" {
		return b.ButtonText.DisplayText
	}
	return b.Label
}

func idForZuckButton(b ZuckInteractiveButtonDTO) string {
	if b.ButtonID != "" {
		return b.ButtonID
	}
	return b.ID
}

// MapZuckButtonsToSendButton — ButtonsMessage (body/text, merchant_url, pix fields).
func MapZuckButtonsToSendButton(d *ZuckButtonsDTO) *send_service.ButtonStruct {
	num := d.ZuckBase.EffectiveRecipient()
	desc := FirstNonEmpty(d.Text, d.Body, d.Message, d.Caption)
	buttons := make([]send_service.Button, 0, len(d.Buttons))
	for _, zb := range d.Buttons {
		typ := normalizeZuckButtonType(zb.Type)
		dt := displayTextForZuckButton(zb)
		id := idForZuckButton(zb)
		ev := send_service.Button{Type: typ, DisplayText: dt, Id: id}
		switch typ {
		case "url":
			ev.URL = FirstNonEmpty(zb.URL, zb.MerchantURL)
			if ev.URL == "" {
				ev.URL = id
			}
		case "call":
			ev.PhoneNumber = zb.Phone
			if ev.PhoneNumber == "" {
				ev.PhoneNumber = id
			}
		case "copy":
			ev.CopyCode = zb.Code
		case "pix":
			ev.Name = zb.MerchantName
			ev.Key = zb.PixKey
			if zb.PixType != "" {
				ev.KeyType = strings.ToLower(zb.PixType)
			}
			ev.Currency = zb.Currency
		case "review_and_pay":
			// não mapeado para Evolution; validador deve bloquear antes do envio
		}
		buttons = append(buttons, ev)
	}
	out := &send_service.ButtonStruct{
		Title:       d.Title,
		Description: desc,
		Footer:      d.Footer,
		Buttons:     buttons,
	}
	applyZuckBaseToButton(out, &d.ZuckBase, num)
	return out
}

func normalizeCarouselButtonType(t string) string {
	u := strings.ToUpper(strings.TrimSpace(t))
	switch u {
	case "QUICK_REPLY", "REPLY", "":
		return "REPLY"
	case "CTA_URL", "URL":
		return "URL"
	case "CTA_CALL", "CALL":
		return "CALL"
	case "CTA_COPY", "COPY":
		return "COPY"
	default:
		if u == "" {
			return "REPLY"
		}
		return u
	}
}

// MapZuckCarouselToSendCarousel — CarouselMessage (Message, Carousel[], MediaType document → erro).
func MapZuckCarouselToSendCarousel(d *ZuckCarouselDTO) (*send_service.CarouselStruct, error) {
	num := d.ZuckBase.EffectiveRecipient()
	body := effectiveCarouselIntro(d)
	cardIn := d.effectiveCarouselCards()
	cards := make([]send_service.CarouselCardStruct, 0, len(cardIn))
	for i, c := range cardIn {
		ctext := effectiveCardText(c)
		if strings.TrimSpace(ctext) == "" {
			return nil, fmt.Errorf("carousel[%d]: text is required", i)
		}
		mt := strings.ToLower(strings.TrimSpace(effectiveCardMediaType(c)))
		murl := effectiveCardMediaURL(c)
		header := send_service.CarouselCardHeaderStruct{Title: " ", Subtitle: ""}
		if murl != "" {
			switch mt {
			case "", "image":
				header.ImageUrl = murl
			case "video":
				header.VideoUrl = murl
			case "document":
				return nil, fmt.Errorf("carousel[%d]: mediaType document não é suportado pelo núcleo Evolution (apenas image e video no carrossel)", i)
			default:
				return nil, fmt.Errorf("carousel[%d]: unsupported mediaType %q (use image, video or omit)", i, mt)
			}
		}
		if header.ImageUrl != "" && header.VideoUrl != "" {
			return nil, fmt.Errorf("carousel[%d]: specify only one of image or video media", i)
		}
		card := send_service.CarouselCardStruct{
			Header: header,
			Body:   send_service.CarouselCardBodyStruct{Text: ctext},
			Footer: FirstNonEmpty(effectiveCardCaption(c), effectiveCardFilename(c)),
		}
		for _, zb := range effectiveCardButtons(c) {
			typ := normalizeCarouselButtonType(effectiveCarouselBtnType(zb))
			dt := effectiveCarouselBtnLabel(zb)
			cb := send_service.CarouselButtonStruct{Type: typ, DisplayText: dt}
			switch typ {
			case "URL":
				cb.Id = FirstNonEmpty(zb.URL, effectiveCarouselBtnID(zb))
			case "CALL":
				cb.Id = FirstNonEmpty(zb.Phone, effectiveCarouselBtnID(zb))
			case "COPY":
				cb.CopyCode = FirstNonEmpty(zb.Code, effectiveCarouselBtnID(zb), zb.IdCap)
				cb.Id = effectiveCarouselBtnID(zb)
			case "REPLY":
				cb.Id = FirstNonEmpty(effectiveCarouselBtnID(zb), zb.Label, zb.LabelCap)
			default:
				cb.Id = effectiveCarouselBtnID(zb)
			}
			card.Buttons = append(card.Buttons, cb)
		}
		cards = append(cards, card)
	}
	out := &send_service.CarouselStruct{
		Body:   body,
		Footer: d.Footer,
		Cards:  cards,
	}
	applyZuckBaseToCarousel(out, &d.ZuckBase, num)
	return out, nil
}
