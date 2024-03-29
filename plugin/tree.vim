if exists('g:loaded_tree')
    finish
endif
let g:loaded_tree = 1
let s:plugin_root = fnamemodify(resolve(expand('<sfile>:p')), ':h:h')

function! s:getBinName()
  if has("win32")
    return "tree-windows-amd64"
  elseif has("mac") && system('arch') == 'arm64'
    return 'tree-darwin-arm64'
  elseif has("mac")
    return 'tree-darwin-amd64'
  elseif has("unix")
    return "tree-linux-amd64"
  endif
endfunction

function! s:StartPlugin(host) abort
  return jobstart(s:plugin_root.'/bin/' . s:getBinName(), {'rpc': v:true})
endfunction

call remote#host#Register('tree', 'x', function('s:StartPlugin'))
call remote#host#RegisterPlugin('tree', '0', [
\ {'type': 'function', 'name': 'Handler_2f4eab2fca750d49cc9039bac724b89b', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'OperatorFunc_2f4eab2fca750d49cc9039bac724b89b', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'TreeClose', 'sync': 0, 'opts': {}},
\ {'type': 'function', 'name': 'TreeFocus', 'sync': 0, 'opts': {}},
\ {'type': 'function', 'name': 'TreeOpen', 'sync': 0, 'opts': {}},
\ {'type': 'function', 'name': 'TreeToggle', 'sync': 0, 'opts': {}},
\ {'type': 'function', 'name': 'TreeToggleFocus', 'sync': 0, 'opts': {}},
\ {'type': 'function', 'name': 'TreeToggleSmart', 'sync': 0, 'opts': {}},
\ {'type': 'function', 'name': 'TreeUnfocus', 'sync': 0, 'opts': {}},
\ ])
