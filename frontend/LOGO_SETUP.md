# 🎨 Configuração de Logos - GoRound

## Status Atual

✅ **Placeholders integrados:**
- Header: Logo com 32x32px (h-8)
- Hero/Page Principal: Logo com 80x80px (h-20)
- Favicon: SVG + PNG fallback

## Como Substituir pelos Logos Reais

### 1. Gerar o Logo via Google Nano Banana

**Prompt sugerido:**
```
Go programming language Gopher mascot, cute and friendly character, 
colorful, professional logo design, isolated on transparent background
```

**Configurações:**
- Resolução: **512x512px ou 1024x1024px**
- Formato: **PNG com transparência** ou **SVG**
- Background: Transparente

### 2. Preparar os Arquivos

**Opção A: Se gerou PNG**

Redimensione para múltiplos tamanhos:

```bash
# Instalar ImageMagick se não tiver
brew install imagemagick  # macOS
# ou
apt-get install imagemagick  # Linux

# Redimensionar
convert go-gopher-original.png -resize 64x64 go-gopher-64.png
convert go-gopher-original.png -resize 128x128 go-gopher-128.png
convert go-gopher-original.png -resize 256x256 go-gopher-256.png
convert go-gopher-original.png -resize 512x512 go-gopher-512.png
```

**Opção B: Se gerou SVG**
- Use como está (escalável infinitamente)
- Salve como `go-gopher.svg`

### 3. Colocar os Arquivos

Coloque todos no diretório:
```
frontend/static/logos/
├── go-gopher-placeholder.svg  ← REMOVA ou substitua
├── go-gopher.svg             ← NOVO (se SVG)
├── go-gopher-64.png          ← NOVO
├── go-gopher-128.png         ← NOVO
├── go-gopher-256.png         ← NOVO
└── go-gopher-512.png         ← NOVO
```

### 4. Atualizar os Arquivos do Projeto

**Em `src/routes/+layout.svelte` (line ~29):**
```svelte
<!-- Trocar de: -->
<img src="/logos/go-gopher-placeholder.svg" alt="GoRound" class="h-8 w-8" />

<!-- Para: -->
<img src="/logos/go-gopher-64.png" alt="GoRound" class="h-8 w-8" />
```

**Em `src/routes/+page.svelte` (line ~40):**
```svelte
<!-- Trocar de: -->
<img src="/logos/go-gopher-placeholder.svg" alt="GoRound Gopher" class="h-20 w-20 mx-auto mb-4 opacity-80" />

<!-- Para: -->
<img src="/logos/go-gopher-128.png" alt="GoRound Gopher" class="h-20 w-20 mx-auto mb-4 opacity-80" />
```

**Em `src/app.html` (line ~5):**
```html
<!-- Trocar de: -->
<link rel="icon" href="%sveltekit.assets%/logos/go-gopher-placeholder.svg" type="image/svg+xml" />
<link rel="alternate icon" href="%sveltekit.assets%/logos/go-gopher-64.png" type="image/png" />

<!-- Para: -->
<link rel="icon" href="%sveltekit.assets%/logos/go-gopher-64.png" type="image/png" />
<!-- Ou se usar SVG: -->
<link rel="icon" href="%sveltekit.assets%/logos/go-gopher.svg" type="image/svg+xml" />
```

### 5. Testar Localmente

```bash
npm run dev
```

Abra http://localhost:5173 e verifique:
- ✅ Logo aparece no header
- ✅ Logo aparece na página principal
- ✅ Favicon está correto na aba do navegador

### 6. Commit das Mudanças

```bash
git add frontend/static/logos/
git commit -m "Add Go Gopher logo assets"
```

## 📐 Tamanhos Usados no Projeto

| Localização | Tamanho CSS | Arquivo | Dimensão Real |
|---|---|---|---|
| Favicon | - | go-gopher-64.png | 64x64px |
| Header | h-8 (32px) | go-gopher-64.png | 64x64px |
| Hero/Page | h-20 (80px) | go-gopher-128.png | 128x128px |

## 🎯 Paleta de Cores Recomendada

Se customizar o logo, use a paleta GoRound:

```
#E63946 → Arena Red (destaque/CTA)
#00ACD7 → Go Blue (identidade)
#4ADE80 → Code Green (sucesso)
#F1FAEE → Frost White (text)
#111111 → Void Black (fundo)
```

## 📝 Checklist

- [ ] Gerar logo via Google Nano Banana
- [ ] Salvar em `/static/logos/`
- [ ] Redimensionar conforme necessário
- [ ] Atualizar `+layout.svelte`
- [ ] Atualizar `+page.svelte`
- [ ] Atualizar `app.html`
- [ ] Testar em desktop
- [ ] Testar em mobile
- [ ] Testar favicon na aba
- [ ] Commit das mudanças
- [ ] Push para repositório

---

**Dúvidas?** Verifique `static/logos/README.md` para mais detalhes.
