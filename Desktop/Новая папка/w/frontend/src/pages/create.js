import { api, getUser, requireAuth } from '../api.js'

if (!requireAuth()) throw new Error('auth')
const user = getUser()
if (user?.is_admin) window.location.href = '/admin.html'

const form = document.getElementById('record-form')
const alertBox = document.getElementById('form-alert')

function hint(name, msg) {
  const el = document.getElementById(`${name}-hint`)
  if (el) el.textContent = msg || ''
}

form.addEventListener('submit', async (e) => {
  e.preventDefault()
  alertBox.classList.add('d-none')
  hint("address", '')
  hint("dish", '')
  hint("delivery_time", '')
  hint("total", '')
  const payload = {}
  payload["address"] = form["address"].value.trim()
  payload["dish"] = form["dish"].value.trim()
  payload["delivery_time"] = form["delivery_time"].value.trim()
  payload["total"] = Number(form["total"].value)
  let ok = true
  if (!payload["address"] || (typeof payload["address"] === 'string' && !payload["address"])) {
    hint("address", "Адрес доставки: обязательно")
    ok = false
  }
  if (!payload["dish"] || (typeof payload["dish"] === 'string' && !payload["dish"])) {
    hint("dish", "Блюдо: обязательно")
    ok = false
  }
  if (!payload["delivery_time"] || (typeof payload["delivery_time"] === 'string' && !payload["delivery_time"])) {
    hint("delivery_time", "Время доставки: обязательно")
    ok = false
  }
  if (!payload["total"] || (typeof payload["total"] === 'string' && !payload["total"])) {
    hint("total", "Сумма, ₽: обязательно")
    ok = false
  }
  if (!/^\d{2}\.\d{2}\.\d{4}$/.test(payload["delivery_time"])) {
    hint("delivery_time", 'Формат ДД.ММ.ГГГГ')
    ok = false
  }
  if (!ok) return
  try {
    await api('/records', { method: 'POST', body: JSON.stringify(payload) })
    alertBox.textContent = 'Запись создана'
    alertBox.className = 'alert alert-success alert-inline'
    alertBox.classList.remove('d-none')
    form.reset()
    setTimeout(() => (window.location.href = '/cabinet.html'), 700)
  } catch (err) {
    alertBox.textContent = err.message
    alertBox.className = 'alert alert-warning alert-inline'
    alertBox.classList.remove('d-none')
  }
})
