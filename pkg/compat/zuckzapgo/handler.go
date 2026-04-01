package zuckzapgo

import (
	"net/http"
	"strings"

	instance_model "github.com/EvolutionAPI/evolution-go/pkg/instance/model"
	send_service "github.com/EvolutionAPI/evolution-go/pkg/sendMessage/service"
	"github.com/gin-gonic/gin"
)

// ZuckCompatHandler exposes ZuckZapGo-compatible routes under /chat/send/*.
type ZuckCompatHandler interface {
	SendTextCompat(c *gin.Context)
	SendImageCompat(c *gin.Context)
	SendAudioCompat(c *gin.Context)
	SendDocumentCompat(c *gin.Context)
	SendVideoCompat(c *gin.Context)
	SendPTVCompat(c *gin.Context)
	SendStickerCompat(c *gin.Context)
	SendLocationCompat(c *gin.Context)
	SendContactCompat(c *gin.Context)
	SendListCompat(c *gin.Context)
	SendButtonsCompat(c *gin.Context)
	SendCarouselCompat(c *gin.Context)
}

type zuckCompatHandler struct {
	svc send_service.SendService
}

func instanceFromContext(c *gin.Context) (*instance_model.Instance, bool) {
	v, ok := c.Get("instance")
	if !ok {
		return nil, false
	}
	inst, ok := v.(*instance_model.Instance)
	return inst, ok
}

func jsonBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func respondSendError(c *gin.Context, err error) {
	c.JSON(MapSendErrorToStatus(err), gin.H{"error": err.Error()})
}

// SendTextCompat godoc
// @Summary      Send text (ZuckZapGo compat)
// @Description  Compatibilidade ZuckZapGo: aceita `phone` ou `number`, `text`/`message`/`body`. Reutiliza o serviço de /send/text.
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckTextDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/text [post]
func (h *zuckCompatHandler) SendTextCompat(c *gin.Context) {
	inst, ok := instanceFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "instance not found"})
		return
	}
	var dto ZuckTextDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	if err := ValidateZuckText(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	msg, err := h.svc.SendText(MapZuckTextToSendText(&dto), inst)
	if err != nil {
		respondSendError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": msg})
}

// SendImageCompat godoc
// @Summary      Send image (ZuckZapGo compat)
// @Description  Aceita `phone`, `url` ou `image`, `caption`. Reutiliza SendMediaUrl com type=image.
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckMediaDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/image [post]
func (h *zuckCompatHandler) SendImageCompat(c *gin.Context) {
	h.sendMediaCompat(c, "image")
}

// SendAudioCompat godoc
// @Summary      Send audio (ZuckZapGo compat)
// @Description  Aceita `phone`, `url` ou `audio`. O núcleo Evolution envia áudio como nota de voz (PTT), como em /send/media.
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckMediaDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/audio [post]
func (h *zuckCompatHandler) SendAudioCompat(c *gin.Context) {
	h.sendMediaCompat(c, "audio")
}

// SendDocumentCompat godoc
// @Summary      Send document (ZuckZapGo compat)
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckMediaDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/document [post]
func (h *zuckCompatHandler) SendDocumentCompat(c *gin.Context) {
	h.sendMediaCompat(c, "document")
}

// SendVideoCompat godoc
// @Summary      Send video (ZuckZapGo compat)
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckMediaDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/video [post]
func (h *zuckCompatHandler) SendVideoCompat(c *gin.Context) {
	h.sendMediaCompat(c, "video")
}

// SendPTVCompat godoc
// @Summary      Send PTV / vídeo circular (ZuckZapGo compat)
// @Description  Usa o mesmo fluxo que /send/media com type=ptv (vídeo MP4). Suportado pelo núcleo Evolution.
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckMediaDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/ptv [post]
func (h *zuckCompatHandler) SendPTVCompat(c *gin.Context) {
	h.sendMediaCompat(c, "ptv")
}

func (h *zuckCompatHandler) sendMediaCompat(c *gin.Context, evolutionType string) {
	inst, ok := instanceFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "instance not found"})
		return
	}
	var dto ZuckMediaDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	if err := ValidateZuckMedia(&dto, evolutionType); err != nil {
		jsonBadRequest(c, err)
		return
	}
	payload := MapZuckMediaDTO(&dto, evolutionType)
	msg, err := h.svc.SendMediaUrl(payload, inst)
	if err != nil {
		respondSendError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": msg})
}

// SendStickerCompat godoc
// @Summary      Send sticker (ZuckZapGo compat)
// @Description  Aceita URL http(s) ou data:image/...;base64,...
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckStickerDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/sticker [post]
func (h *zuckCompatHandler) SendStickerCompat(c *gin.Context) {
	inst, ok := instanceFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "instance not found"})
		return
	}
	var dto ZuckStickerDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	if err := ValidateZuckSticker(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	msg, err := h.svc.SendSticker(MapZuckStickerToSendSticker(&dto), inst)
	if err != nil {
		respondSendError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": msg})
}

// SendLocationCompat godoc
// @Summary      Send location (ZuckZapGo compat)
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckLocationDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/location [post]
func (h *zuckCompatHandler) SendLocationCompat(c *gin.Context) {
	inst, ok := instanceFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "instance not found"})
		return
	}
	var dto ZuckLocationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	if err := ValidateZuckLocation(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	msg, err := h.svc.SendLocation(MapZuckLocationToSendLocation(&dto), inst)
	if err != nil {
		respondSendError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": msg})
}

// SendContactCompat godoc
// @Summary      Send contact (ZuckZapGo compat)
// @Description  Aceita campos flat ou objeto `vcard` no formato Evolution.
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckContactDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/contact [post]
func (h *zuckCompatHandler) SendContactCompat(c *gin.Context) {
	inst, ok := instanceFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "instance not found"})
		return
	}
	var dto ZuckContactDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	if err := ValidateZuckContact(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	msg, err := h.svc.SendContact(MapZuckContactToSendContact(&dto), inst)
	if err != nil {
		respondSendError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": msg})
}

// SendListCompat godoc
// @Summary      Send list (ZuckZapGo compat)
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckListDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/list [post]
func (h *zuckCompatHandler) SendListCompat(c *gin.Context) {
	inst, ok := instanceFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "instance not found"})
		return
	}
	var dto ZuckListDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	if err := ValidateZuckList(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	msg, err := h.svc.SendList(MapZuckListToSendList(&dto), inst)
	if err != nil {
		respondSendError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": msg})
}

// SendButtonsCompat godoc
// @Summary      Send buttons (ZuckZapGo compat)
// @Description  Tipos Zuck: quick_reply, cta_url, cta_call, cta_copy mapeados para o fluxo de /send/button.
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckButtonsDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/buttons [post]
func (h *zuckCompatHandler) SendButtonsCompat(c *gin.Context) {
	inst, ok := instanceFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "instance not found"})
		return
	}
	var dto ZuckButtonsDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	if err := ValidateZuckButtons(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	msg, err := h.svc.SendButton(MapZuckButtonsToSendButton(&dto), inst)
	if err != nil {
		respondSendError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": msg})
}

// SendCarouselCompat godoc
// @Summary      Send carousel (ZuckZapGo compat)
// @Description  Expõe SendCarousel do núcleo (não havia rota nativa /send/carousel).
// @Tags         ZuckZapGo compat
// @Accept       json
// @Produce      json
// @Param        body  body      ZuckCarouselDTO  true  "Payload"
// @Success      200   {object}  gin.H
// @Failure      400   {object}  gin.H
// @Failure      500   {object}  gin.H
// @Router       /chat/send/carousel [post]
func (h *zuckCompatHandler) SendCarouselCompat(c *gin.Context) {
	inst, ok := instanceFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "instance not found"})
		return
	}
	var dto ZuckCarouselDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	if err := ValidateZuckCarousel(&dto); err != nil {
		jsonBadRequest(c, err)
		return
	}
	payload, err := MapZuckCarouselToSendCarousel(&dto)
	if err != nil {
		jsonBadRequest(c, err)
		return
	}
	msg, err := h.svc.SendCarousel(payload, inst)
	if err != nil {
		respondSendError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": msg})
}

// NewZuckCompatHandler builds the ZuckZapGo compatibility handler.
func NewZuckCompatHandler(svc send_service.SendService) ZuckCompatHandler {
	return &zuckCompatHandler{svc: svc}
}

// MapSendErrorToStatus maps known Evolution send errors to HTTP status (optional use).
func MapSendErrorToStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	s := err.Error()
	if strings.Contains(s, "not registered on WhatsApp") {
		return http.StatusBadRequest
	}
	if strings.Contains(s, "could not parse") || strings.Contains(s, "Invalid") {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
