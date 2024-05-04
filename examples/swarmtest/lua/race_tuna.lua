function race_tuna(str)
  str = string.gsub(str, "motor sport", "tuna hunters")
  str = string.gsub(str, "drivers", "hunters")
  str = string.gsub(str, "driver", "hunter")
  str = string.gsub(str, "racing", "tuna hunting")
  str = string.gsub(str, "race", "tuna")
  return str
end
