syn match TreeIcon     /\(^\(  \)*. \)\@<=[^ ]/
syn match TreeDirIcon  /[ﱮ▸▾•]/ containedin=TreeIcon
syn match TreeFileIcon /[•]/    containedin=TreeIcon

syn match TreeName     /\(^\(  \)*. [ﱮ▸▾•] \)\@<=.*$/
syn match TreeDirName  /\(^\(  \)*. [ﱮ▸▾•] \)\@<=.*$/
syn match TreeFileName /\(^\(  \)*. [•] \)\@<=.*$/
syn match TreeDirSlash #/# containedin=TreeName,TreeDirName

syn match TreeStatus            /\(^\(  \)*\)\@<=[^ ]\([^ ] \)\@=/
syn match TreeStatusChanged     /\(^\(  \)*\)\@<=◎/  containedin=TreeStatus
syn match TreeStatusAdded       /\(^\(  \)*\)\@<=⦿/  containedin=TreeStatus
syn match TreeStatusConcflicted /\(^\(  \)*\)\@<=◉/  containedin=TreeStatus

" Default theme
highlight default link TreeFileIcon  Normal
highlight default link TreeDirIcon   Directory
highlight default link TreeDirSlash  Comment
highlight default link TreeDirName   Directory

highlight default link TreeStatus            Comment
highlight default link TreeStatusChanged     TreeStatus
highlight default link TreeStatusAdded       TreeStatus
highlight default link TreeStatusConcflicted Error

