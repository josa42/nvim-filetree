call plug#begin()
Plug '~/github/nvim-filetree'
call plug#end()

nnoremap <silent>b   :call TreeToggle()<CR>

nnoremap <silent><Tab>      :tabnext<CR>
nnoremap <silent><S-Tab>    :tabprevious<CR>
nnoremap <silent>t          :tabe<CR>
nnoremap <silent>x          :tabc<CR>

set mouse=a

set splitbelow
set splitright
set number
