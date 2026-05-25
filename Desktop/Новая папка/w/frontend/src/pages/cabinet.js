import { api, clearSession, getUser, requireAuth } from '../api.js'
import { initSlider } from '../slider.js'
import { STATUS_NEW } from '../theme-config.js'

if (!requireAuth()) throw new Error('auth')
const user = getUser()
if (user?.is_admin) window.location.href = '/admin.html'
document.getElementById('user-greeting').textContent = `Здравствуйте, ${user.full_name}`
document.getElementById('logout-btn').onclick = () => { clearSession(); location.href = '/index.html' }
initSlider(document.getElementById('hero-slider'))

const listEl = document.getElementById('records-list')
const reviewModal = new bootstrap.Modal(document.getElementById('reviewModal'))
const reviewForm = document.getElementById('review-form')

function canReview(status) {
  return status !== STATUS_NEW
}

async function load() {
  const rows = await api('/records/my')
  if (!rows.length) { listEl.innerHTML = '<p class="text-muted small">Пока пусто</p>'; return }
  listEl.innerHTML = rows.map((r) => {
    const info = [`${r.address ?? ''}`, `${r.dish ?? ''}`, `${r.delivery_time ?? ''}`, `${r.total ?? ''}`].join(' · ')
    let actions = ''
    actions += `<span class="status-badge">${r.status || ''}</span>`
    if (canReview(r.status) && !r.Review) actions += `<button class="btn btn-sm btn-outline-primary mt-2 d-block" data-review="${r.ID}">Отзыв</button>`
    else if (r.Review) actions += `<p class="small text-muted mt-2 mb-0">${r.Review.text} (${r.Review.rating}/5)</p>`
    return `<article class="record-item"><div class="d-flex justify-content-between align-items-start"><strong>#${r.ID}</strong><div>${actions}</div></div><div class="small mt-1">${info}</div></article>`
  }).join('')
  listEl.querySelectorAll('[data-review]').forEach((btn) => {
    btn.onclick = () => { reviewForm.record_id.value = btn.dataset.review; reviewModal.show() }
  })
}
reviewForm.onsubmit = async (e) => {
  e.preventDefault()
  try {
    await api('/reviews', { method: 'POST', body: JSON.stringify({
      record_id: Number(reviewForm.record_id.value),
      text: reviewForm.text.value,
      rating: Number(reviewForm.rating.value),
    })})
    reviewModal.hide()
    load()
  } catch (err) { document.getElementById('review-hint').textContent = err.message }
}

load().catch(() => { listEl.innerHTML = '<p class="text-danger small">Ошибка загрузки</p>' })
