# GoRound Logos

## Estrutura de Logos

Este diretório contém os logos do GoRound em diferentes formatos e tamanhos.

### Logos Necessários

| Arquivo | Tamanho | Uso | Formato |
|---------|---------|-----|---------|
| `go-gopher.svg` | Vetorial | Mascote/Branding | SVG |
| `go-gopher-64.png` | 64x64px | Favicon | PNG |
| `go-gopher-128.png` | 128x128px | Page Hero | PNG |
| `go-gopher-256.png` | 256x256px | Social/Marketing | PNG |
| `go-gopher-512.png` | 512x512px | Alta Resolução | PNG |

### Como Gerar via Google Nano Banana

1. Acesse: https://nanobanana.google.com/ (ou similar)
2. **Prompt sugerido:**
   ```
   Go programming language Gopher mascot, cute, colorful, 
   friendly character, logo style, professional design
   ```
3. Gere em **1024x1024px** ou **512x512px**
4. Faça download

### Como Usar

#### Para SVG:
- Coloque o arquivo como `go-gopher.svg`
- Use em `+layout.svelte` para logo escalável

#### Para PNG:
- Redimensione conforme necessário usando ImageMagick:
  ```bash
  convert go-gopher-original.png -resize 64x64 go-gopher-64.png
  convert go-gopher-original.png -resize 128x128 go-gopher-128.png
  convert go-gopher-original.png -resize 256x256 go-gopher-256.png
  ```

### Integração no Projeto

#### 1. Favicon (favicon.ico)
```html
<!-- Em src/app.html -->
<link rel="icon" href="/logos/go-gopher-64.png" />
```

#### 2. Logo no Header
```svelte
<!-- Em src/routes/+layout.svelte -->
<img src="/logos/go-gopher-64.png" alt="GoRound" class="h-8" />
```

#### 3. Logo na Página Principal
```svelte
<!-- Em src/routes/+page.svelte -->
<img src="/logos/go-gopher-128.png" alt="GoRound" class="mb-4" />
```

### Guia de Cores (Paleta GoRound)

Quando gerar o logo, considere usar:
- **Arena Red**: #E63946 (destaque)
- **Go Blue**: #00ACD7 (identidade Go)
- **Code Green**: #4ADE80 (sucesso)
- **Frost White**: #F1FAEE (texto)
- **Void Black**: #111111 (fundo)

### Próximos Passos

- [ ] Gerar logo via Google Nano Banana
- [ ] Salvar como `go-gopher.svg` (ou `.png`)
- [ ] Criar variações de tamanho (64, 128, 256, 512)
- [ ] Integrar no `+layout.svelte` (header)
- [ ] Integrar no `+page.svelte` (hero)
- [ ] Atualizar `app.html` com favicon
- [ ] Testar em diferentes resoluções de tela
