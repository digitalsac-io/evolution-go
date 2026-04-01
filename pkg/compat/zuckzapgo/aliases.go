package zuckzapgo

import send_service "github.com/EvolutionAPI/evolution-go/pkg/sendMessage/service"

// ZuckContextInfo mirrors OpenAPI definitions/ContextInfo (replies, mentions).
type ZuckContextInfo struct {
	StanzaID    string `json:"stanzaID"`
	Participant string `json:"participant"`
	MentionedJID []string `json:"mentionedJID"`
	MentionAll  bool   `json:"mentionAll"`
	IsForwarded bool   `json:"isForwarded"`
}

// EffectiveRecipient resolves phone/number per Zuck (Phone) e Evolution (phone, number).
func (b *ZuckBase) EffectiveRecipient() string {
	return FirstNonEmpty(b.Number, b.Phone, b.PhonePascal)
}

// EffectiveMessageID resolves id / Id (OpenAPI PascalCase).
func (b *ZuckBase) EffectiveMessageID() string {
	return FirstNonEmpty(b.ID, b.IdPascal)
}

// EffectiveDelay: Evolution delay (ms) ou Presence do Zuck (ms de digitação).
func (b *ZuckBase) EffectiveDelay() int32 {
	if b.Delay != 0 {
		return b.Delay
	}
	if b.Presence != 0 {
		return b.Presence
	}
	return b.PresencePascal
}

func (b *ZuckBase) activeContext() *ZuckContextInfo {
	if b.ContextInfo != nil {
		return b.ContextInfo
	}
	return b.ContextInfoSnake
}

// EffectiveQuotedAndMentions merges quoted (Evolution) com ContextInfo (Zuck OpenAPI).
func (b *ZuckBase) EffectiveQuotedAndMentions() (send_service.QuotedStruct, []string, bool) {
	q := b.Quoted
	mentions := b.MentionedJid
	ma := b.MentionAll

	if ctx := b.activeContext(); ctx != nil {
		if ctx.StanzaID != "" {
			q.MessageID = ctx.StanzaID
		}
		if ctx.Participant != "" {
			q.Participant = ctx.Participant
		}
		if len(ctx.MentionedJID) > 0 {
			mentions = ctx.MentionedJID
		}
		if ctx.MentionAll {
			ma = true
		}
	}
	return q, mentions, ma
}

// ApplyZuckBaseToTextStruct preenche campos comuns a partir do ZuckBase efetivo.
func ApplyZuckBaseToTextStruct(t *send_service.TextStruct, b *ZuckBase, number string) {
	q, mentions, ma := b.EffectiveQuotedAndMentions()
	t.Number = number
	t.Id = b.EffectiveMessageID()
	t.Delay = b.EffectiveDelay()
	t.FormatJid = b.FormatJid
	t.MentionedJID = mentions
	t.MentionAll = ma
	t.Quoted = q
}

func applyZuckBaseToMedia(m *send_service.MediaStruct, b *ZuckBase, number string) {
	q, mentions, ma := b.EffectiveQuotedAndMentions()
	m.Number = number
	m.Id = b.EffectiveMessageID()
	m.Delay = b.EffectiveDelay()
	m.FormatJid = b.FormatJid
	m.MentionedJID = mentions
	m.MentionAll = ma
	m.Quoted = q
}

func applyZuckBaseToSticker(m *send_service.StickerStruct, b *ZuckBase, number string) {
	q, mentions, ma := b.EffectiveQuotedAndMentions()
	m.Number = number
	m.Id = b.EffectiveMessageID()
	m.Delay = b.EffectiveDelay()
	m.FormatJid = b.FormatJid
	m.MentionedJID = mentions
	m.MentionAll = ma
	m.Quoted = q
}

func applyZuckBaseToLocation(m *send_service.LocationStruct, b *ZuckBase, number string) {
	q, mentions, ma := b.EffectiveQuotedAndMentions()
	m.Number = number
	m.Id = b.EffectiveMessageID()
	m.Delay = b.EffectiveDelay()
	m.FormatJid = b.FormatJid
	m.MentionedJID = mentions
	m.MentionAll = ma
	m.Quoted = q
}

func applyZuckBaseToContact(m *send_service.ContactStruct, b *ZuckBase, number string) {
	q, mentions, ma := b.EffectiveQuotedAndMentions()
	m.Number = number
	m.Id = b.EffectiveMessageID()
	m.Delay = b.EffectiveDelay()
	m.FormatJid = b.FormatJid
	m.MentionedJID = mentions
	m.MentionAll = ma
	m.Quoted = q
}

func applyZuckBaseToList(m *send_service.ListStruct, b *ZuckBase, number string) {
	q, mentions, ma := b.EffectiveQuotedAndMentions()
	m.Number = number
	m.Id = b.EffectiveMessageID()
	m.Delay = b.EffectiveDelay()
	m.FormatJid = b.FormatJid
	m.MentionedJID = mentions
	m.MentionAll = ma
	m.Quoted = q
}

func applyZuckBaseToButton(m *send_service.ButtonStruct, b *ZuckBase, number string) {
	q, mentions, ma := b.EffectiveQuotedAndMentions()
	m.Number = number
	m.Delay = b.EffectiveDelay()
	m.FormatJid = b.FormatJid
	m.MentionedJID = mentions
	m.MentionAll = ma
	m.Quoted = q
}

func applyZuckBaseToCarousel(m *send_service.CarouselStruct, b *ZuckBase, number string) {
	q, _, _ := b.EffectiveQuotedAndMentions()
	m.Number = number
	m.Delay = b.EffectiveDelay()
	m.FormatJid = b.FormatJid
	m.Quoted = q
}

// --- MessageList / Carousel helpers (PascalCase OpenAPI) ---

func (d *ZuckListDTO) effectiveSections() []ZuckListSectionDTO {
	if len(d.SectionsCap) > 0 {
		return d.SectionsCap
	}
	return d.Sections
}

func effectiveSectionTitle(s ZuckListSectionDTO) string {
	return FirstNonEmpty(s.Title, s.TitleCap)
}

func effectiveSectionRows(s ZuckListSectionDTO) []ZuckListRowDTO {
	if len(s.RowsCap) > 0 {
		return s.RowsCap
	}
	return s.Rows
}

func effectiveRowTitle(r ZuckListRowDTO) string {
	return FirstNonEmpty(r.Title, r.TitleCap)
}

func effectiveRowDescription(r ZuckListRowDTO) string {
	return FirstNonEmpty(r.Description, r.DescriptionCap)
}

func effectiveRowID(r ZuckListRowDTO) string {
	return FirstNonEmpty(r.RowID, r.RowIDCap, r.ID)
}

func (d *ZuckCarouselDTO) effectiveCarouselCards() []ZuckCarouselCardDTO {
	if len(d.CarouselCap) > 0 {
		return d.CarouselCap
	}
	return d.Carousel
}

func effectiveCarouselIntro(d *ZuckCarouselDTO) string {
	return FirstNonEmpty(d.Message, d.MessageCap, d.Text, d.Body)
}

func effectiveCardText(c ZuckCarouselCardDTO) string {
	return FirstNonEmpty(c.Text, c.TextCap)
}

func effectiveCardMediaURL(c ZuckCarouselCardDTO) string {
	return FirstNonEmpty(c.MediaURL, c.MediaURLCap)
}

func effectiveCardMediaType(c ZuckCarouselCardDTO) string {
	return FirstNonEmpty(c.MediaType, c.MediaTypeCap)
}

func effectiveCardFilename(c ZuckCarouselCardDTO) string {
	return FirstNonEmpty(c.Filename, c.FileNameCap)
}

func effectiveCardCaption(c ZuckCarouselCardDTO) string {
	return FirstNonEmpty(c.Caption, c.CaptionCap)
}

func effectiveCardButtons(c ZuckCarouselCardDTO) []ZuckCarouselCardButtonDTO {
	if len(c.ButtonsCap) > 0 {
		return c.ButtonsCap
	}
	return c.Buttons
}

func effectiveCarouselBtnLabel(b ZuckCarouselCardButtonDTO) string {
	return FirstNonEmpty(b.DisplayText, b.Label, b.LabelCap)
}

func effectiveCarouselBtnType(b ZuckCarouselCardButtonDTO) string {
	return FirstNonEmpty(b.Type, b.TypeCap)
}

func effectiveCarouselBtnID(b ZuckCarouselCardButtonDTO) string {
	return FirstNonEmpty(b.ID, b.IdCap)
}
