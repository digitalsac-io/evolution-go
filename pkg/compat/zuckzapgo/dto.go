package zuckzapgo

import send_service "github.com/EvolutionAPI/evolution-go/pkg/sendMessage/service"

// ZuckBase: campos Evolution (lowerCamel) + aliases do OpenAPI ZuckZapGo (PascalCase / snake).
type ZuckBase struct {
	Phone   string `json:"phone"`
	Number  string `json:"number"`
	ID      string `json:"id"`
	Delay   int32  `json:"delay"`
	FormatJid *bool `json:"formatJid"`

	MentionedJid []string `json:"mentionedJid"`
	MentionAll   bool     `json:"mentionAll"`
	Quoted       send_service.QuotedStruct `json:"quoted"`

	// OpenAPI Zuck (MessageText, MessageImage, …)
	PhonePascal    string `json:"Phone"`
	IdPascal       string `json:"Id"`
	Presence       int32  `json:"presence"`
	PresencePascal int32  `json:"Presence"`
	ContextInfo    *ZuckContextInfo `json:"ContextInfo"`
	ContextInfoSnake *ZuckContextInfo `json:"context_info"`
}

// ZuckTextDTO — "#/definitions/MessageText" (+ aliases Evolution).
type ZuckTextDTO struct {
	ZuckBase
	Text        string `json:"text"`
	Message     string `json:"message"`
	Body        string `json:"body"`
	BodyPascal  string `json:"Body"`
}

// ZuckMediaURLRef — header de mídia em ButtonsMessage (image/video/document).
type ZuckMediaURLRef struct {
	URL string `json:"url"`
}

// ZuckMediaDTO — MessageImage, MessageAudio, MessageVideo, MessagePTV, MessageDocument.
type ZuckMediaDTO struct {
	ZuckBase
	URL      string `json:"url"`
	Image    string `json:"image"`
	ImageCap string `json:"Image"`
	Audio    string `json:"audio"`
	AudioCap string `json:"Audio"`
	Video    string `json:"video"`
	VideoCap string `json:"Video"`
	Document string `json:"document"`
	DocCap   string `json:"Document"`

	Caption    string `json:"caption"`
	CaptionCap string `json:"Caption"`

	Filename    string `json:"filename"`
	FileNameCap string `json:"FileName"`

	Mimetype    string `json:"mimetype"`
	MimeTypeCap string `json:"MimeType"`

	PTT    *bool `json:"ptt"`
	PTTCap *bool `json:"PTT"`

	IsPtv  bool   `json:"isPtv"`
	Base64 string `json:"base64"`
}

// ZuckStickerDTO — "#/definitions/MessageSticker"
type ZuckStickerDTO struct {
	ZuckBase
	Sticker     string `json:"sticker"`
	StickerCap  string `json:"Sticker"`
	URL         string `json:"url"`
	MimeTypeCap string `json:"MimeType"`
}

// ZuckLocationDTO — "#/definitions/MessageLocation" (+ address para Evolution).
type ZuckLocationDTO struct {
	ZuckBase
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	LatitudeCap  float64 `json:"Latitude"`
	LongitudeCap float64 `json:"Longitude"`
	Name         string  `json:"name"`
	NameCap      string  `json:"Name"`
	Address      string  `json:"address"`
}

// ZuckContactDTO — "#/definitions/MessageContact" (Phone, Name, Vcard string) + formato Evolution.
type ZuckContactDTO struct {
	ZuckBase
	// Destino: já coberto por ZuckBase (phone / Phone / number)

	NameCap  string `json:"Name"`
	VcardRaw string `json:"Vcard"`

	FullName     string `json:"fullName"`
	Organization string `json:"organization"`
	ContactPhone string `json:"contactPhone"`
	Vcard        *struct {
		FullName     string `json:"fullName"`
		Organization string `json:"organization"`
		Phone        string `json:"phone"`
	} `json:"vcard"`
}

// ZuckListRowDTO — rows em MessageList (rowId / RowId).
type ZuckListRowDTO struct {
	Title       string `json:"title"`
	TitleCap    string `json:"Title"`
	Description string `json:"description"`
	DescriptionCap string `json:"Description"`
	RowID       string `json:"rowId"`
	RowIDCap    string `json:"RowId"`
	ID          string `json:"id"`
}

// ZuckListSectionDTO — sections / Sections.
type ZuckListSectionDTO struct {
	Title    string            `json:"title"`
	TitleCap string            `json:"Title"`
	Rows     []ZuckListRowDTO `json:"rows"`
	RowsCap  []ZuckListRowDTO `json:"Rows"`
}

// ZuckListDTO — "#/definitions/MessageList"
type ZuckListDTO struct {
	ZuckBase
	Title       string `json:"title"`
	TitleCap    string `json:"Title"`
	Description string `json:"description"`
	Text        string `json:"text"`
	TextCap     string `json:"Text"`
	Body        string `json:"body"`
	ButtonText  string `json:"buttonText"`
	ButtonTextCap string `json:"ButtonText"`
	Footer      string `json:"footer"`
	FooterCap   string `json:"Footer"`
	FooterText  string `json:"footerText"`
	Sections    []ZuckListSectionDTO `json:"sections"`
	SectionsCap []ZuckListSectionDTO `json:"Sections"`
}

// ZuckButtonTextDTO
type ZuckButtonTextDTO struct {
	DisplayText string `json:"displayText"`
}

// ZuckInteractiveButtonDTO — item de buttons[] em ButtonsMessage (+ merchant_url, pix_*, etc.).
type ZuckInteractiveButtonDTO struct {
	ButtonID    string            `json:"buttonId"`
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	ButtonText  ZuckButtonTextDTO `json:"buttonText"`
	Label       string            `json:"label"`
	URL         string            `json:"url"`
	MerchantURL string            `json:"merchant_url"`
	Phone       string            `json:"phone"`
	Code        string            `json:"code"`

	PixKey       string `json:"pix_key"`
	MerchantName string `json:"merchant_name"`
	PixType      string `json:"pix_type"`
	Currency     string `json:"currency"`
	TotalValue   int    `json:"total_value"`
	TotalOffset  int    `json:"total_offset"`
	ReferenceID  string `json:"reference_id"`
}

// ZuckButtonsDTO — "#/definitions/ButtonsMessage"
type ZuckButtonsDTO struct {
	ZuckBase
	Title   string `json:"title"`
	Text    string `json:"text"`
	Body    string `json:"body"`
	Message string `json:"message"`
	Footer  string `json:"footer"`
	Caption string `json:"caption"`

	Image    *ZuckMediaURLRef `json:"image"`
	Video    *ZuckMediaURLRef `json:"video"`
	Document *ZuckMediaURLRef `json:"document"`

	Buttons []ZuckInteractiveButtonDTO `json:"buttons"`
}

// ZuckCarouselCardButtonDTO — "#/definitions/CarouselMessage" card buttons (Id / Label / Type).
type ZuckCarouselCardButtonDTO struct {
	ID          string `json:"id"`
	IdCap       string `json:"Id"`
	Label       string `json:"label"`
	LabelCap    string `json:"Label"`
	DisplayText string `json:"displayText"`
	Type        string `json:"type"`
	TypeCap     string `json:"Type"`
	URL         string `json:"url"`
	Phone       string `json:"phone"`
	Code        string `json:"code"`
}

// ZuckCarouselCardDTO
type ZuckCarouselCardDTO struct {
	Text        string                      `json:"text"`
	TextCap     string                      `json:"Text"`
	MediaURL    string                      `json:"mediaUrl"`
	MediaURLCap string                      `json:"MediaUrl"`
	MediaType   string                      `json:"mediaType"`
	MediaTypeCap string                     `json:"MediaType"`
	Filename    string                      `json:"filename"`
	FileNameCap string                      `json:"Filename"`
	Caption     string                      `json:"caption"`
	CaptionCap  string                      `json:"Caption"`
	Buttons     []ZuckCarouselCardButtonDTO `json:"buttons"`
	ButtonsCap  []ZuckCarouselCardButtonDTO `json:"Buttons"`
}

// ZuckCarouselDTO — "#/definitions/CarouselMessage"
type ZuckCarouselDTO struct {
	ZuckBase
	Message       string                `json:"message"`
	MessageCap    string                `json:"Message"`
	Text          string                `json:"text"`
	Body          string                `json:"body"`
	Footer        string                `json:"footer"`
	Carousel      []ZuckCarouselCardDTO `json:"carousel"`
	CarouselCap   []ZuckCarouselCardDTO `json:"Carousel"`
}
