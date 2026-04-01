package zuckzapgo

import (
	"fmt"
	"strings"
)

func validateRecipientBase(b *ZuckBase) error {
	if strings.TrimSpace(b.EffectiveRecipient()) == "" {
		return fmt.Errorf("phone, Phone or number is required")
	}
	return nil
}

func isAllowedMediaPayload(s string) bool {
	low := strings.ToLower(strings.TrimSpace(s))
	return strings.HasPrefix(low, "http://") ||
		strings.HasPrefix(low, "https://") ||
		strings.HasPrefix(low, "data:")
}

// ValidateZuckText — MessageText (Body obrigatório no OpenAPI).
func ValidateZuckText(d *ZuckTextDTO) error {
	if err := validateRecipientBase(&d.ZuckBase); err != nil {
		return err
	}
	if strings.TrimSpace(FirstNonEmpty(d.Text, d.Message, d.Body, d.BodyPascal)) == "" {
		return fmt.Errorf("text, message, body or Body is required")
	}
	return nil
}

// ValidateZuckMedia — URL, data: ou campo base64 (quando preenche URL no mapper).
func ValidateZuckMedia(d *ZuckMediaDTO, evolutionType string) error {
	if err := validateRecipientBase(&d.ZuckBase); err != nil {
		return err
	}
	url := pickMediaURL(d, evolutionType)
	if strings.TrimSpace(d.Base64) != "" && strings.TrimSpace(url) == "" {
		url = d.Base64
	}
	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("media URL or data URL is required (url, image/Image, audio/Audio, etc.)")
	}
	if !isAllowedMediaPayload(url) {
		return fmt.Errorf("media must be an http(s) URL or a data: URL (OpenAPI MessageImage/Audio/...)")
	}
	return nil
}

// ValidateZuckSticker
func ValidateZuckSticker(d *ZuckStickerDTO) error {
	if err := validateRecipientBase(&d.ZuckBase); err != nil {
		return err
	}
	st := FirstNonEmpty(d.Sticker, d.StickerCap, d.URL)
	if strings.TrimSpace(st) == "" {
		return fmt.Errorf("sticker, Sticker or url is required")
	}
	if !isAllowedMediaPayload(st) {
		return fmt.Errorf("sticker must be an http(s) URL or a data: URL")
	}
	return nil
}

// ValidateZuckLocation
func ValidateZuckLocation(d *ZuckLocationDTO) error {
	if err := validateRecipientBase(&d.ZuckBase); err != nil {
		return err
	}
	lat := d.Latitude
	if lat == 0 && d.LatitudeCap != 0 {
		lat = d.LatitudeCap
	}
	lng := d.Longitude
	if lng == 0 && d.LongitudeCap != 0 {
		lng = d.LongitudeCap
	}
	if lat == 0 && lng == 0 {
		return fmt.Errorf("latitude and longitude are required (latitude/Latitude, longitude/Longitude)")
	}
	name := FirstNonEmpty(d.Name, d.NameCap)
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name or Name is required")
	}
	return nil
}

// ValidateZuckContact — MessageContact ou vcard Evolution.
func ValidateZuckContact(d *ZuckContactDTO) error {
	if err := validateRecipientBase(&d.ZuckBase); err != nil {
		return err
	}
	vc := MapZuckContactToSendContact(d)
	if strings.TrimSpace(vc.Vcard.Phone) == "" {
		return fmt.Errorf("contact phone is required (Vcard TEL, vcard.phone, contactPhone)")
	}
	if strings.TrimSpace(vc.Vcard.FullName) == "" {
		return fmt.Errorf("contact full name is required (Name, Vcard FN, fullName, vcard.fullName)")
	}
	return nil
}

// ValidateZuckList — MessageList (Text, ButtonText, Sections obrigatórios; Title opcional).
func ValidateZuckList(d *ZuckListDTO) error {
	if err := validateRecipientBase(&d.ZuckBase); err != nil {
		return err
	}
	desc := FirstNonEmpty(d.Description, d.Text, d.TextCap, d.Body)
	if strings.TrimSpace(desc) == "" {
		return fmt.Errorf("text, Text, description or body is required (OpenAPI MessageList.Text)")
	}
	btn := FirstNonEmpty(d.ButtonText, d.ButtonTextCap)
	if strings.TrimSpace(btn) == "" {
		return fmt.Errorf("buttonText or ButtonText is required")
	}
	footer := FirstNonEmpty(d.FooterText, d.Footer, d.FooterCap)
	if strings.TrimSpace(footer) == "" {
		return fmt.Errorf("footer, Footer or footerText is required")
	}
	secs := d.effectiveSections()
	if len(secs) == 0 {
		return fmt.Errorf("sections or Sections must not be empty")
	}
	for i, sec := range secs {
		rows := effectiveSectionRows(sec)
		if len(rows) == 0 {
			return fmt.Errorf("sections[%d]: rows or Rows must not be empty", i)
		}
		if strings.TrimSpace(effectiveSectionTitle(sec)) == "" {
			return fmt.Errorf("sections[%d]: title or Title is required", i)
		}
		for j, row := range rows {
			if strings.TrimSpace(effectiveRowTitle(row)) == "" {
				return fmt.Errorf("sections[%d].rows[%d]: title or Title is required", i, j)
			}
			if strings.TrimSpace(effectiveRowID(row)) == "" {
				return fmt.Errorf("sections[%d].rows[%d]: rowId, RowId or id is required", i, j)
			}
		}
	}
	return nil
}

func hasButtonMediaHeader(d *ZuckButtonsDTO) bool {
	if d.Image != nil && strings.TrimSpace(d.Image.URL) != "" {
		return true
	}
	if d.Video != nil && strings.TrimSpace(d.Video.URL) != "" {
		return true
	}
	if d.Document != nil && strings.TrimSpace(d.Document.URL) != "" {
		return true
	}
	return false
}

// ValidateZuckButtons — ButtonsMessage (OpenAPI: phone + buttons; title/footer opcionais).
func ValidateZuckButtons(d *ZuckButtonsDTO) error {
	if err := validateRecipientBase(&d.ZuckBase); err != nil {
		return err
	}
	if hasButtonMediaHeader(d) {
		return fmt.Errorf("compat: image/video/document header em botões (OpenAPI ButtonsMessage) não é suportado nesta camada; use sem header de mídia ou evolua o SendButton no núcleo")
	}
	body := FirstNonEmpty(d.Text, d.Body, d.Message, d.Caption)
	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("text, body, message or caption is required when there is no media header")
	}
	if len(d.Buttons) == 0 {
		return fmt.Errorf("buttons must not be empty")
	}
	for i, b := range d.Buttons {
		typ := normalizeZuckButtonType(b.Type)
		dt := displayTextForZuckButton(b)
		if strings.TrimSpace(dt) == "" {
			return fmt.Errorf("buttons[%d]: display text is required (buttonText.displayText or label)", i)
		}
		switch typ {
		case "reply":
			if strings.TrimSpace(idForZuckButton(b)) == "" {
				return fmt.Errorf("buttons[%d]: buttonId or id is required for quick_reply", i)
			}
		case "url":
			url := FirstNonEmpty(b.URL, b.MerchantURL)
			if url == "" {
				url = idForZuckButton(b)
			}
			if strings.TrimSpace(url) == "" {
				return fmt.Errorf("buttons[%d]: url, merchant_url or buttonId is required for cta_url", i)
			}
		case "call":
			ph := b.Phone
			if ph == "" {
				ph = idForZuckButton(b)
			}
			if strings.TrimSpace(ph) == "" {
				return fmt.Errorf("buttons[%d]: phone is required for cta_call", i)
			}
		case "copy":
			if strings.TrimSpace(b.Code) == "" {
				return fmt.Errorf("buttons[%d]: code is required for cta_copy", i)
			}
		case "pix":
			// Evolution valida campos PIX ao enviar
		case "review_and_pay":
			return fmt.Errorf("buttons[%d]: type review_and_pay não é suportado na camada compat Evolution", i)
		default:
			return fmt.Errorf("buttons[%d]: unsupported type %q", i, b.Type)
		}
	}
	return nil
}

// ValidateZuckCarousel — CarouselMessage (Message, Carousel; cada card: Text + Buttons).
func ValidateZuckCarousel(d *ZuckCarouselDTO) error {
	if err := validateRecipientBase(&d.ZuckBase); err != nil {
		return err
	}
	cards := d.effectiveCarouselCards()
	if len(cards) == 0 {
		return fmt.Errorf("carousel or Carousel must contain at least one card")
	}
	if strings.TrimSpace(effectiveCarouselIntro(d)) == "" {
		return fmt.Errorf("message, Message, text or body is required (intro text before cards)")
	}
	for i, c := range cards {
		if strings.TrimSpace(effectiveCardText(c)) == "" {
			return fmt.Errorf("carousel[%d]: text or Text is required", i)
		}
		murl := effectiveCardMediaURL(c)
		if murl != "" {
			mt := strings.ToLower(strings.TrimSpace(effectiveCardMediaType(c)))
			if mt == "document" {
				return fmt.Errorf("carousel[%d]: mediaType document não é suportado pelo núcleo Evolution no carrossel", i)
			}
			if mt != "" && mt != "image" && mt != "video" {
				return fmt.Errorf("carousel[%d]: mediaType must be image or video when MediaUrl is set", i)
			}
		}
		btns := effectiveCardButtons(c)
		if len(btns) == 0 {
			return fmt.Errorf("carousel[%d]: buttons or Buttons must not be empty (OpenAPI CarouselMessage)", i)
		}
		for j, btn := range btns {
			typ := normalizeCarouselButtonType(effectiveCarouselBtnType(btn))
			dt := effectiveCarouselBtnLabel(btn)
			if strings.TrimSpace(dt) == "" {
				return fmt.Errorf("carousel[%d].buttons[%d]: label, Label or displayText is required", i, j)
			}
			switch typ {
			case "URL":
				if strings.TrimSpace(FirstNonEmpty(btn.URL, effectiveCarouselBtnID(btn))) == "" {
					return fmt.Errorf("carousel[%d].buttons[%d]: url or Id is required for URL button", i, j)
				}
			case "CALL":
				if strings.TrimSpace(FirstNonEmpty(btn.Phone, effectiveCarouselBtnID(btn))) == "" {
					return fmt.Errorf("carousel[%d].buttons[%d]: phone or Id is required for CALL button", i, j)
				}
			case "COPY":
				if strings.TrimSpace(FirstNonEmpty(btn.Code, effectiveCarouselBtnID(btn), btn.IdCap)) == "" {
					return fmt.Errorf("carousel[%d].buttons[%d]: code or Id is required for COPY button", i, j)
				}
			case "REPLY":
				if strings.TrimSpace(FirstNonEmpty(effectiveCarouselBtnID(btn), btn.Label, btn.LabelCap)) == "" {
					return fmt.Errorf("carousel[%d].buttons[%d]: id, Id or label is required for reply button", i, j)
				}
			default:
				return fmt.Errorf("carousel[%d].buttons[%d]: unsupported type %q", i, j, effectiveCarouselBtnType(btn))
			}
		}
	}
	return nil
}
